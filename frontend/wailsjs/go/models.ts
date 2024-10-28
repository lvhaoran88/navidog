export namespace types {
	
	export class JSResp {
	    success: boolean;
	    message: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new JSResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = source["data"];
	    }
	}
	export class MysqlConnection {
	    host: string;
	    port: string;
	    username: string;
	    password: string;
	    database?: string;
	
	    static createFrom(source: any = {}) {
	        return new MysqlConnection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.database = source["database"];
	    }
	}

}

