import { useRequest, useUpdateEffect } from "ahooks";
import {
  Button,
  Divider,
  Form,
  Input,
  InputNumber,
  message,
  Modal,
} from "antd";
import { forwardRef, useImperativeHandle, useRef, useState } from "react";

import { TestConnection } from "../../wailsjs/go/services/mysqlService";
import { SaveConnection } from "../../wailsjs/go/services/connectionService";
import { FormInstance } from "antd/lib";
import { types } from "../../wailsjs/go/models";

const ConnectionForm = forwardRef((props, ref) => {
  const [form] = Form.useForm();

  const { loading: testLoading, runAsync: testConnection } = useRequest(
    () => {
      return form.validateFields().then((values) => {
        return TestConnection(values.config);
      });
    },
    {
      manual: true,
      refreshDeps: [form],
    }
  );

  const { loading: saveLoading, runAsync: saveConnection } = useRequest(
    () => {
      return form.validateFields().then((values) => {
        return SaveConnection(values);
      });
    },
    {
      manual: true,
      refreshDeps: [form],
    }
  );

  useImperativeHandle(
    ref,
    () => ({
      form,
      testConnection,
      saveConnection,
      testLoading,
      saveLoading,
    }),
    [form, testConnection, saveConnection, testLoading, saveLoading]
  );

  return (
    <Form
      form={form}
      labelCol={{ span: 4 }}
      initialValues={{
        name: "localhost",
        config: {
          host: "localhost",
          port: 3306,
          username: "",
          password: "",
        },
      }}
    >
      <Form.Item label="名称" name="name">
        <Input placeholder="请输入名称" />
      </Form.Item>
      <Divider />
      <Form.Item label="地址" name={["config", "host"]}>
        <Input placeholder="请输入地址" />
      </Form.Item>
      <Form.Item label="端口" name={["config", "port"]}>
        <InputNumber style={{ width: "100%" }} />
      </Form.Item>
      <Form.Item label="用户名" name={["config", "username"]}>
        <Input placeholder="请输入用户名" />
      </Form.Item>
      <Form.Item label="密码" name={["config", "password"]}>
        <Input.Password placeholder="请输入密码" />
      </Form.Item>
    </Form>
  );
});

interface ConnectionFormRef {
  form: FormInstance;
  testConnection: () => Promise<types.JSResp>;
  saveConnection: () => Promise<types.JSResp>;
  testLoading: boolean;
  saveLoading: boolean;
}

const CreateConnection = () => {
  const [open, setOpen] = useState(false);
  const connectionFormRef = useRef<ConnectionFormRef>();

  return (
    <>
      <Button type="primary" onClick={() => setOpen(true)}>
        新建链接
      </Button>
      <Modal
        title="新建链接"
        open={open}
        okText="保存"
        onOk={() => {
          connectionFormRef.current?.saveConnection().then((res) => {
            if (!res.success) {
              Modal.error({
                title: "保存失败",
                content: res.message,
              });
            } else {
              message.success("保存成功");
              connectionFormRef.current?.form.resetFields();
              setOpen(false);
            }
            return res;
          });
        }}
        okButtonProps={{
          loading: connectionFormRef.current?.saveLoading,
          disabled: connectionFormRef.current?.testLoading,
        }}
        onCancel={() => setOpen(false)}
        cancelButtonProps={{
          loading:
            connectionFormRef.current?.saveLoading ||
            connectionFormRef.current?.testLoading,
        }}
        footer={(_, { OkBtn, CancelBtn }) => {
          return (
            <>
              <Button
                type="primary"
                ghost
                loading={connectionFormRef.current?.testLoading}
                disabled={connectionFormRef.current?.saveLoading}
                onClick={() => {
                  connectionFormRef.current?.testConnection().then((res) => {
                    if (!res.success) {
                      Modal.error({
                        title: "测试连接失败",
                        content: res.message,
                      });
                    } else {
                      Modal.success({
                        title: "测试连接成功",
                      });
                    }
                    return res;
                  });
                }}
              >
                测试连接
              </Button>
              <CancelBtn />
              <OkBtn />
            </>
          );
        }}
      >
        <ConnectionForm ref={connectionFormRef} />
      </Modal>
    </>
  );
};

export default CreateConnection;
