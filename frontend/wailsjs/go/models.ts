export namespace dictionary {
	
	export class Category {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Category(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class Word {
	    id: number;
	    english: string;
	    afaanOromo: string;
	    partOfSpeech: string;
	    exampleEn: string;
	    exampleOm: string;
	    pronunciation: string;
	    categories?: Category[];
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Word(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.english = source["english"];
	        this.afaanOromo = source["afaanOromo"];
	        this.partOfSpeech = source["partOfSpeech"];
	        this.exampleEn = source["exampleEn"];
	        this.exampleOm = source["exampleOm"];
	        this.pronunciation = source["pronunciation"];
	        this.categories = this.convertValues(source["categories"], Category);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class WordInput {
	    english: string;
	    afaanOromo: string;
	    partOfSpeech: string;
	    exampleEn: string;
	    exampleOm: string;
	    pronunciation: string;
	    categoryIds: number[];
	
	    static createFrom(source: any = {}) {
	        return new WordInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.english = source["english"];
	        this.afaanOromo = source["afaanOromo"];
	        this.partOfSpeech = source["partOfSpeech"];
	        this.exampleEn = source["exampleEn"];
	        this.exampleOm = source["exampleOm"];
	        this.pronunciation = source["pronunciation"];
	        this.categoryIds = source["categoryIds"];
	    }
	}

}

