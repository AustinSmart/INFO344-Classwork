"use strict";
const defaultResponse = "Sorry, I'm not sure what you are asking. Try asking with different phrasing";

class MongoInterface {
    /** 
    * @param {mongodb.Collection} messagesCollection 
    * @param {mongodb.Collection} channelsCollection 
    **/
    constructor(messagesCollection, channelsCollection) {
        this.messagesCollection = messagesCollection;
        this.channelsCollection = channelsCollection;
    }
   
    async hasNotPosted(channel) {
     	return "You want who hasnt posted to a channel";
    }

    async mostPost(channel) {
	    return "You want who posted most in a chanel";
    }

    async members(channel) {
        return "You want who is in a channel";
    }

    async lastPostInChannel(user, channel) {
		return "You want your last post in a channel";
    }

    async lastPost(user) {
	    return "You want your last post";
    }

    async totalMessagesInChannelOnDate(user, channel, date) {
        return "You want your total messages in a channel at a date";
    }

    async totalMessagesInChannel(user, channel) {
        var CircularJSON = require('circular-json');
        var channelFromDB = await this.channelsCollection.find(
            {$or: [{name: new RegExp(channel, "i"), members: user.id}, {name: new RegExp(channel, "i"), private: false}]}
            ).toArray();
        if(!channelFromDB[0]) {
            return "Sorry, I can't find that channel or you do not have access to it.";
        }
        console.log(JSON.stringify(user));
        console.log(CircularJSON.stringify(channelFromDB));
        var allMsgs = await this.messagesCollection.find({creatorid: user.id, channelid: channelFromDB[0]._id}).toArray();
        return `You have posted ${allMsgs.length} messages in the "${channelFromDB[0].name}" channel`;
    }

   async totalMessages(user) {
        var allMsgs = await this.messagesCollection.find({creatorid: user.id}).toArray();
        if(!allMsgs) {
            return defaultResponse;
        }
        return `You have posted ${allMsgs.length} messages`;
    }
}

module.exports = MongoInterface;