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
         var channelFromDB = await this.channelsCollection.find(
            {$or: [{name: new RegExp(channel, "i"), members: user.id}, {name: new RegExp(channel, "i"), private: false}]}
            ).toArray();
        if(!channelFromDB[0]) {
            return "Sorry, I can't find that channel or you do not have access to it.";
        }
		 var msg = await this.messagesCollection.find({creatorid: user.id, channelid: channelFromDB[0]._id}).sort({createdat: -1}).limit(1).toArray();
        if(!msg[0]) {
            return `Sorry, I could not find your latest message in ${channel}`;
        }
	    return `Your last post in "${channelFromDB[0].name}" was ${msg[0].createdat}`;
    }

    async lastPost(user) {
        var msg = await this.messagesCollection.find({creatorid: user.id}).sort({createdat: -1}).limit(1).toArray();
        if(!msg[0]) {
            return "Sorry, I could not find your latest message";
        }
	    return `Your last post was ${msg[0].createdat}`;
    }

    async totalMessagesInChannelOnDate(user, channel, date) {
        console.log("DATE: " + new Date(date));
        var channelFromDB = await this.channelsCollection.find(
            {$or: [{name: new RegExp(channel, "i"), members: user.id}, {name: new RegExp(channel, "i"), private: false}]}
            ).toArray();
        if(!channelFromDB[0]) {
            return "Sorry, I can't find that channel or you do not have access to it.";
        }
        var allMsgs = await this.messagesCollection.find(
            {creatorid: user.id, 
            channelid: channelFromDB[0]._id, 
            createdat: new Date(date).toISOString()})
            .toArray();
        return `You posted ${allMsgs.length} messages in the "${channelFromDB[0].name}" channel on ${new Date(date)}`;
    }

    async totalMessagesInChannel(user, channel) {
        var channelFromDB = await this.channelsCollection.find(
            {$or: [{name: new RegExp(channel, "i"), members: user.id}, {name: new RegExp(channel, "i"), private: false}]}
            ).toArray();
        if(!channelFromDB[0]) {
            return "Sorry, I can't find that channel or you do not have access to it.";
        }
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