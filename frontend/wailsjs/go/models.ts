export namespace main {
	
	export class Account {
	    email: string;
	    source: string;
	    is_service_account: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.email = source["email"];
	        this.source = source["source"];
	        this.is_service_account = source["is_service_account"];
	    }
	}
	export class AuthImportRequest {
	    users: any[];
	    hashAlgorithm: string;
	    signerKey: string;
	    saltSeparator: string;
	    rounds: number;
	    memoryCost: number;
	
	    static createFrom(source: any = {}) {
	        return new AuthImportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.users = source["users"];
	        this.hashAlgorithm = source["hashAlgorithm"];
	        this.signerKey = source["signerKey"];
	        this.saltSeparator = source["saltSeparator"];
	        this.rounds = source["rounds"];
	        this.memoryCost = source["memoryCost"];
	    }
	}
	export class AuthImportResult {
	    success: boolean;
	    importedCount: number;
	    failedCount: number;
	    errors?: any[];
	
	    static createFrom(source: any = {}) {
	        return new AuthImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.importedCount = source["importedCount"];
	        this.failedCount = source["failedCount"];
	        this.errors = source["errors"];
	    }
	}
	export class AuthUser {
	    uid: string;
	    email: string;
	    emailVerified: boolean;
	    displayName: string;
	    photoUrl: string;
	    disabled: boolean;
	    createdAt: string;
	    lastLoginAt: string;
	    customClaims: Record<string, any>;
	    customAttributes?: string;
	    passwordHash?: string;
	    salt?: string;
	    providers?: string[];
	    providerUserInfo?: any[];
	
	    static createFrom(source: any = {}) {
	        return new AuthUser(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uid = source["uid"];
	        this.email = source["email"];
	        this.emailVerified = source["emailVerified"];
	        this.displayName = source["displayName"];
	        this.photoUrl = source["photoUrl"];
	        this.disabled = source["disabled"];
	        this.createdAt = source["createdAt"];
	        this.lastLoginAt = source["lastLoginAt"];
	        this.customClaims = source["customClaims"];
	        this.customAttributes = source["customAttributes"];
	        this.passwordHash = source["passwordHash"];
	        this.salt = source["salt"];
	        this.providers = source["providers"];
	        this.providerUserInfo = source["providerUserInfo"];
	    }
	}
	export class AuthUsersResponse {
	    users: AuthUser[];
	    nextPageToken?: string;
	
	    static createFrom(source: any = {}) {
	        return new AuthUsersResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.users = this.convertValues(source["users"], AuthUser);
	        this.nextPageToken = source["nextPageToken"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExportAuthResponse {
	    users: AuthUser[];
	    count: number;
	    hashParameters: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new ExportAuthResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.users = this.convertValues(source["users"], AuthUser);
	        this.count = source["count"];
	        this.hashParameters = source["hashParameters"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Project {
	    projectId: string;
	    projectNumber: string;
	    displayName: string;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectId = source["projectId"];
	        this.projectNumber = source["projectNumber"];
	        this.displayName = source["displayName"];
	    }
	}

}

