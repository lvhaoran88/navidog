import { useRequest } from "ahooks";
import { message, Select } from "antd";
import { SelectProps } from "antd/lib";
import { types } from "../../wailsjs/go/models";
import { ListDatabases } from "../../wailsjs/go/services/mysqlService";
import { useMemo } from "react";

interface DatabaseSelectProps extends SelectProps<string> {
  connection?: types.Connection;
}

const DatabaseSelect: React.FC<DatabaseSelectProps> = ({
  connection,
  ...props
}) => {
  const { loading, data = [] } = useRequest<string[], []>(
    () => {
      if (!connection) return Promise.resolve([]);
      return ListDatabases(connection.config).then((res) => {
        if (!res.success) {
          message.error(res.message);
          return [];
        }
        return res.data;
      });
    },
    { refreshDeps: [connection] }
  );
  const options = useMemo(() => {
    return data.map((item) => {
      return {
        label: item,
        value: item,
      };
    });
  }, [data]);
  return <Select options={options} {...props} />;
};

export default DatabaseSelect;
