import {
  Button,
  Card,
  Col,
  ConfigProvider,
  Descriptions,
  Form,
  Input,
  message,
  Modal,
  Progress,
  Row,
} from "antd";
import React, { useMemo, useState } from "react";
import DatabaseSelect from "../components/DatabaseSelect";
import { types } from "../../wailsjs/go/models";
import TableTree from "../components/TableTree";
import ConnectionSelect from "../components/ConnectionSlect";
import { TransferData as xxx } from "../../wailsjs/go/services/mysqlService";
import { useCounter, useMount, useRequest } from "ahooks";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { useMysqlVersion } from "../hooks/mysql";

const MysqlVersionInfo = ({
  connection,
  version,
}: {
  connection?: types.Connection;
  version?: string;
}) => {
  if (!connection || !version) return null;
  return (
    <ConfigProvider
      theme={{
        components: {
          Descriptions: {
            itemPaddingBottom: 0,
            itemPaddingEnd: 0,
          },
        },
      }}
    >
      <Descriptions
        items={[
          {
            label: "名称",
            key: "name",
            children: connection.name,
          },
          {
            label: "主机",
            key: "host",
            children: connection.config.host,
          },
          {
            label: "端口",
            key: "port",
            children: connection.config.port,
          },
          {
            label: "用户",
            key: "user",
            children: connection.config.username,
          },
          {
            label: "版本",
            key: "version",
            children: version,
          },
        ]}
        column={1}
      />
    </ConfigProvider>
  );
};

const TransferDataForm = ({ onFinish }: { onFinish: () => void }) => {
  const [sourceConnection, setSourceConnection] = useState<types.Connection>();
  const { data: sourceVersion } = useMysqlVersion<string | undefined, []>(
    sourceConnection?.config
  );
  const [targetConnection, setTargetConnection] = useState<types.Connection>();
  const { data: targetVersion } = useMysqlVersion<string | undefined, []>(
    targetConnection?.config
  );

  const [sourceDatabase, setSourceDatabase] = useState<string>();
  const [targetDatabase, setTargetDatabase] = useState<string>();

  const [tables, setTables] = useState<React.Key[]>([]);

  const mergeSourceConnection = useMemo(() => {
    if (!sourceConnection || !sourceDatabase) {
      return undefined;
    }
    return {
      ...sourceConnection,
      config: {
        ...sourceConnection.config,
        database: sourceDatabase,
      },
    } as types.Connection;
  }, [sourceConnection, sourceDatabase]);

  const mergeTargetConnection = useMemo(() => {
    if (!targetConnection || !targetDatabase) {
      return undefined;
    }
    return {
      ...targetConnection,
      config: {
        ...targetConnection.config,
        database: targetDatabase,
      },
    } as types.Connection;
  }, [targetConnection, targetDatabase]);

  const { loading, runAsync: handlerTransferData } = useRequest(
    () => {
      if (!mergeSourceConnection || !mergeTargetConnection || !tables) {
        console.log("mergeSourceConnection: ", mergeSourceConnection);
        console.log("mergeTargetConnection: ", mergeTargetConnection);
        console.log("tables: ", tables);
        return Promise.reject("请选择源数据、目标数据、数据库和表");
      }
      return xxx(
        mergeSourceConnection.config,
        mergeTargetConnection.config,
        tables as string[]
      ).then((res) => {
        console.log("transfer data res: ", res);
        return res;
      });
    },
    {
      manual: true,
      refreshDeps: [mergeSourceConnection, mergeTargetConnection, tables],
    }
  );

  const [percent, setPercent] = useState(0);
  const [log, setLog] = useState<string[]>([]);
  useMount(() => {
    EventsOn("transferDataProgress", (index: number, count: number) => {
      console.log("transfer data progress: ", index, count);
      setPercent(Number(((index / count) * 100).toFixed(2)));
    });
    EventsOn("transferDataLog", (data: string) => {
      console.log("transfer data log: ", data);
      setLog((p) => [...p, data]);
    });
  });

  const [step, { inc, dec }] = useCounter(1, { min: 1, max: 5 });

  // 操作按钮
  const handlerButton = useMemo(() => {
    const finish1 =
      sourceConnection && sourceDatabase && targetConnection && targetDatabase;
    const finish2 = finish1 && tables;
    if (step === 1) {
      return (
        <Button
          type="primary"
          disabled={!finish1}
          onClick={() => {
            inc();
          }}
        >
          下一步
        </Button>
      );
    } else if (step === 2) {
      return (
        <>
          <Button
            onClick={() => {
              dec();
            }}
          >
            上一步
          </Button>
          <Button
            type="primary"
            disabled={!finish2}
            onClick={() => {
              inc();
            }}
          >
            下一步
          </Button>
        </>
      );
    } else if (step === 3) {
      return (
        <>
          <Button
            onClick={() => {
              dec();
            }}
          >
            上一步
          </Button>
          <Button
            onClick={() => {
              inc();
              handlerTransferData().then((res) => {
                if (!res.success) {
                  message.error(res.message);
                } else {
                  inc();
                }
              });
            }}
            type="primary"
          >
            开始
          </Button>
        </>
      );
    } else if (step === 4) {
      return (
        <Button
          onClick={() => {
            console.log("transfer data finish");
          }}
          type="primary"
          danger
        >
          终止
        </Button>
      );
    } else {
      return (
        <Button onClick={onFinish} type="primary">
          完成
        </Button>
      );
    }
  }, [
    step,
    handlerTransferData,
    sourceConnection,
    sourceDatabase,
    targetConnection,
    targetDatabase,
    tables,
    onFinish,
  ]);

  return (
    <Form component={false}>
      <Form.Item hidden={step !== 1}>
        <Row gutter={[24, 0]}>
          <Col span={12}>
            <Card title="源">
              <Form.Item label="链接" layout="vertical">
                <ConnectionSelect onChange={setSourceConnection} />
              </Form.Item>
              <Form.Item label="数据库" layout="vertical">
                <DatabaseSelect
                  connection={sourceConnection}
                  onChange={setSourceDatabase}
                />
              </Form.Item>
              <Form.Item label="信息" layout="vertical">
                <MysqlVersionInfo
                  connection={sourceConnection}
                  version={sourceVersion}
                />
              </Form.Item>
            </Card>
          </Col>
          <Col span={12}>
            <Card title="目标">
              <Form.Item label="连接" layout="vertical">
                <ConnectionSelect onChange={setTargetConnection} />
              </Form.Item>
              <Form.Item label="数据库" layout="vertical">
                <DatabaseSelect
                  connection={targetConnection}
                  onChange={setTargetDatabase}
                />
              </Form.Item>
              <Form.Item label="信息" layout="vertical">
                <MysqlVersionInfo
                  connection={targetConnection}
                  version={targetVersion}
                />
              </Form.Item>
            </Card>
          </Col>
        </Row>
      </Form.Item>
      <Form.Item label="选择表" layout="vertical" hidden={step !== 2}>
        <TableTree
          connection={mergeSourceConnection}
          checkedKeys={tables}
          // @ts-ignore
          onCheck={(_, { checked, node }) => {
            setTables((p) => {
              if (checked) {
                return [...p, node.key];
              } else {
                return p.filter((i) => i !== node.key);
              }
            });
          }}
        />
      </Form.Item>
      <Form.Item
        label="日志"
        layout="vertical"
        hidden={step === 1 || step === 2}
      >
        <Input.TextArea value={log.join("\r\n")} rows={8} />
        <Progress percent={percent} status="active" />
      </Form.Item>
      {handlerButton}
    </Form>
  );
};

const TransferData = () => {
  const [open, setOpen] = useState(false);
  return (
    <>
      <Button type="primary" onClick={() => setOpen(true)}>
        数据传输
      </Button>
      <Modal
        width={740}
        title="数据传输"
        open={open}
        onCancel={() => setOpen(false)}
        footer={null}
        destroyOnClose
      >
        <TransferDataForm onFinish={() => setOpen(false)} />
      </Modal>
    </>
  );
};

export default TransferData;
