import { useDeepCompareEffect, useMount, useRequest } from "ahooks";
import { Input, Modal, Spin, Tree, TreeNodeProps } from "antd";
import { types } from "../../wailsjs/go/models";
import { ListTables } from "../../wailsjs/go/services/mysqlService";
import { useMemo, useState } from "react";
import { DataNode } from "antd/es/tree";
import { TreeProps } from "antd/lib";

interface TableTreeProps extends TreeProps {
  connection?: types.Connection;
}

const TableTree: React.FC<TableTreeProps> = ({
  connection,
  checkedKeys,
  onCheck,
}) => {
  const {
    loading,
    data = [],
    runAsync,
  } = useRequest<string[], []>(
    () => {
      if (!connection) return Promise.resolve([]);
      return ListTables(connection.config).then((res) => {
        if (!res.success) {
          Modal.error({
            content: res.message,
          });
          return [];
        }
        return res.data;
      });
    },
    {
      manual: true,
      refreshDeps: [connection],
    }
  );

  useMount(() => {
    runAsync();
  });

  useDeepCompareEffect(() => {
    console.log("connection", connection);
    runAsync();
  }, [connection]);

  const [searchValue, setSearchValue] = useState("");

  const treeData: DataNode[] = useMemo(() => {
    if (!connection) return [];

    const tableNodeData: DataNode[] = [];
    if (!searchValue) {
      data.forEach((table) => {
        return tableNodeData.push({
          key: table,
          title: table,
        });
      });
    } else {
      data.forEach((table) => {
        if (table.includes(searchValue)) {
          tableNodeData.push({
            key: table,
            title: table,
          });
        }
      });
    }

    return [
      {
        checkable: false,
        key: `connection-name-${connection.name}'}`,
        title: `${connection.name}(${connection.config.database})`,
        children: [
          {
            checkable: false,
            key: `database-name-${connection.config.database}`,

            title: `表(${(checkedKeys as React.Key[])?.length || 0}/${
              data.length
            })`,
            children: tableNodeData,
          },
        ],
      },
    ];
  }, [connection, checkedKeys, data, searchValue]);

  return (
    <Spin spinning={loading}>
      <Input.Search
        allowClear
        placeholder="搜索"
        onChange={(e) => setSearchValue(e.target.value)}
        onSearch={setSearchValue}
      />
      <Tree height={400} checkable treeData={treeData} onCheck={onCheck} />
    </Spin>
  );
};

export default TableTree;
