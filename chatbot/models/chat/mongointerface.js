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
        var channelFromDB = await this.channelsCollection.findOne({name: channel});
        console.log(CircularJSON.stringify(channelFromDB) + channelFromDB.channelid);
        var allMsgs = await this.messagesCollection.find({creatorid:user.id}, {channel:channelFromDB.channelid}).toArray();
        return `You have posted ${allMsgs.length} messages in the ${channel} channel`;
    }

   async totalMessages(user) {
        var allMsgs = await this.messagesCollection.find({creatorid:user.id}).toArray();
        return `You have posted ${allMsgs.length} messages`;
    }
}

module.exports = MongoInterface;