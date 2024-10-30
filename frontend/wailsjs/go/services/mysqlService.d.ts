// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {types} from '../models';
import {context} from '../models';

export function ListDatabases(arg1:types.MysqlConnection):Promise<types.JSResp>;

export function ListTables(arg1:types.MysqlConnection):Promise<types.JSResp>;

export function Startup(arg1:context.Context):Promise<void>;

export function TestConnection(arg1:types.MysqlConnection):Promise<types.JSResp>;

export function TransferData(arg1:types.MysqlConnection,arg2:types.MysqlConnection,arg3:Array<string>):Promise<types.JSResp>;

export function Version(arg1:types.MysqlConnection):Promise<types.JSResp>;
