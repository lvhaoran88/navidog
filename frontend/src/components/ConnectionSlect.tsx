import { useRequest } from "ahooks";
import { Modal, Select } from "antd";
import { ListConnections } from "../../wailsjs/go/services/connectionService";
import { useCallback, useMemo } from "react";
import { types } from "../../wailsjs/go/models";
import { SelectProps } from "antd/lib";

interface ConnectionSelect extends SelectProps {
  value?: types.Connection;
  onChange?: (value?: types.Connection) => void;
}

const ConnectionSelect: React.FC<ConnectionSelect> = ({
  value: restValue,
  onChange: restOnChange,
  ...props
}) => {
  const { loading, data } = useRequest<types.Connection[], []>(
    () => {
      return ListConnections().then((res) => {
        if (!res.success) {
          Modal.error({
            title: "Error",
            content: res.message,
          });
          return [];
        }
        return res.data;
      });
    },
    {
      refreshDeps: [],
    }
  );

  const options = useMemo(() => {
    return data?.map((item) => ({
      label: item.name,
      value: JSON.stringify(item),
    }));
  }, [data]);

  const value = useMemo(() => {
    if (!restValue) return undefined;
    return JSON.stringify(restValue);
  }, [restValue]);

  const onChange = useCallback(
    (value: string | undefined) => {
      if (!value || !restOnChange) return;
      restOnChange?.(JSON.parse(value));
    },
    [restOnChange]
  );

  return (
    <Select
      options={options}
      loading={loading}
      value={value}
      onChange={onChange}
      {...props}
    />
  );
};

export default ConnectionSelect;
