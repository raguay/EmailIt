export namespace main {
	
	export class FileInfo {
	    Dir: string;
	    Name: string;
	    Extension: string;
	    IsDir: boolean;
	    Size: number;
	    Modtime: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Dir = source["Dir"];
	        this.Name = source["Name"];
	        this.Extension = source["Extension"];
	        this.IsDir = source["IsDir"];
	        this.Size = source["Size"];
	        this.Modtime = source["Modtime"];
	    }
	}
	export class FileParts {
	    Dir: string;
	    Name: string;
	    Extension: string;
	
	    static createFrom(source: any = {}) {
	        return new FileParts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Dir = source["Dir"];
	        this.Name = source["Name"];
	        this.Extension = source["Extension"];
	    }
	}
	export class GitHubRepos {
	    name: string;
	    url: string;
	    stars: number;
	    owner: string;
	    id: number;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new GitHubRepos(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.url = source["url"];
	        this.stars = source["stars"];
	        this.owner = source["owner"];
	        this.id = source["id"];
	        this.description = source["description"];
	    }
	}

}

