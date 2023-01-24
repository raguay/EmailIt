export namespace main {
	
	
	
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

