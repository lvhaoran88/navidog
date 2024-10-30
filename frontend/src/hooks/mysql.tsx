import { useRequest } from "ahooks";
import { types } from "../../wailsjs/go/models";
import { Version as MysqlVersion } from "../../wailsjs/go/services/mysqlService";
import { Options } from "ahooks/lib/useRequest/src/types";

export function useMysqlVersion<TData, TParams extends any[]>(
  config: types.MysqlConnection | undefined,
  option?: Options<TData, TParams>
) {
  return useRequest<TData, TParams>(
    () => {
      if (!config) return Promise.resolve(undefined);
      return MysqlVersion(config).then((res) => {
        if (!res.success) {
          console.error(res.message);
          return undefined;
        }
        return res.data;
      });
    },
    {
      refreshDeps: [config],
      ...(option || {}),
    }
  );
}
