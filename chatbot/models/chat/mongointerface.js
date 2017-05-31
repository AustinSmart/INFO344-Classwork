"use strict";

class MongoInterface {
    /** 
    * @param {mongodb.Collection} messagesCollection 
    * @param {mongodb.Collection} channelsCollection 
    **/
    constructor(messagesCollection, channelsCollection) {
        this.messagesCollection = messagesCollection;
        this.channelsCollection = channelsCollection;
    }
   
    hasNotPosted(channel) {
     	return "You want who hasnt posted to a channel";
    }

    mostPost(channel) {
	    return "You want who posted most in a chanel";
    }

    members(channel) {
        return "You want who is in a channel";
    }

    lastPostInChannel(user, channel) {
		return "You want your last post in a channel";
    }

    lastPost(user) {
	    return "You want your last post";
    }

    totalMessageInChannelOnDate(user, channel, date) {
        return "You want your total messages in a channel at a date";
    }

    totalMessageInChannel(user, channel) {
        return "You want your total messages in a channel";
    }

    totalMessages(user) {
        var allMsgs = this.messagesCollection.find().toArray();
        return `You have posted ${allMsgs.length} messages`;
    }

}

module.exports = MongoInterface;