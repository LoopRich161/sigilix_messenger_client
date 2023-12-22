export namespace data {
	
	export class Message {
	    message_id: number;
	    chat_id: number;
	    sender_id: number;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.message_id = source["message_id"];
	        this.chat_id = source["chat_id"];
	        this.sender_id = source["sender_id"];
	        this.content = source["content"];
	    }
	}
	export class Chat {
	    chat_id: number;
	    other_user_id: number;
	    last_message_id: number;
	    am_i_initiator: boolean;
	    accepted: boolean;
	    title: string;
	    messages: Message[];
	
	    static createFrom(source: any = {}) {
	        return new Chat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.chat_id = source["chat_id"];
	        this.other_user_id = source["other_user_id"];
	        this.last_message_id = source["last_message_id"];
	        this.am_i_initiator = source["am_i_initiator"];
	        this.accepted = source["accepted"];
	        this.title = source["title"];
	        this.messages = this.convertValues(source["messages"], Message);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

}

export namespace messenger_client {
	
	export class WebNotificationWithTypeInfo {
	    notification: any;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new WebNotificationWithTypeInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.notification = source["notification"];
	        this.type = source["type"];
	    }
	}

}

