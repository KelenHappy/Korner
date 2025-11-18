export namespace main {
	
	export class AppSettings {
	    apiProvider: string;
	    apiKey: string;
	    apiEndpoint: string;
	    floatingIcon: string;
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.apiProvider = source["apiProvider"];
	        this.apiKey = source["apiKey"];
	        this.apiEndpoint = source["apiEndpoint"];
	        this.floatingIcon = source["floatingIcon"];
	    }
	}

}

