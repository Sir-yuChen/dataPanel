export namespace code {
	
	export class App {
	    Handler: any;
	
	    static createFrom(source: any = {}) {
	        return new App(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Handler = source["Handler"];
	    }
	}

}

export namespace configModel {
	
	export class System {
	    applicationName: string;
	    env: string;
	    addr: number;
	    "db-type": string;
	    "use-multipoint": boolean;
	
	    static createFrom(source: any = {}) {
	        return new System(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.applicationName = source["applicationName"];
	        this.env = source["env"];
	        this.addr = source["addr"];
	        this["db-type"] = source["db-type"];
	        this["use-multipoint"] = source["use-multipoint"];
	    }
	}

}

export namespace exposed {
	
	export class HelloWails {
	
	
	    static createFrom(source: any = {}) {
	        return new HelloWails(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace options {
	
	export class SecondInstanceData {
	    Args: string[];
	    WorkingDirectory: string;
	
	    static createFrom(source: any = {}) {
	        return new SecondInstanceData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Args = source["Args"];
	        this.WorkingDirectory = source["WorkingDirectory"];
	    }
	}

}

