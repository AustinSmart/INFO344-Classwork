"use strict";
const defaultResponse = "Sorry, I'm not sure what you are asking. Try asking with different phrasing";
const CirularJSON = require('circular-json');
const moment = require('moment');

class MongoInterface {
    /** 
    * @param {mongodb.Collection} usersCollection 
    * @param {mongodb.Collection} messagesCollection 
    * @param {mongodb.Collection} channelsCollection 
    **/
    constructor(usersCollection, messagesCollection, channelsCollection) {
        this.usersCollection = usersCollection;
        this.messagesCollection = messagesCollection;
        this.channelsCollection = channelsCollection;
    }
   
    async hasNotPosted(user, channel) {
        var channelFromDB = await this.channelsCollection.find(
            {$or: [{name: new RegExp(channel, "i"), members: user.id}, {name: new RegExp(channel, "i"), private: false}]}
            ).toArray();
        if(!channelFromDB[0]) {
            return "Sorry, I can't find that channel or you do not have access to it.";
        }

        var allMsgs = await this.messagesCollection.find({channelid: channelFromDB[0]._id}).toArray();
        var creatorIDs = allMsgs.map(function(m) {return m.creatorid;});
        if(channelFromDB[0].private === false) {
            var users = await this.usersCollection.find().toArray();
        } else {
            var users = new Array();
            for (var userID of channelFromDB[0].members) {
                var user = await this.usersCollection.findOne({_id: userID});
                if(user) {
                    users.push(user);
                }
            }
        }
        var userIDs = users.map(function(u) {return u._id;});
        var hasNotPostedIDs = userIDs.filter(function(u){return this.indexOf(u)<0;},creatorIDs);
        var hasNotPostedUsers = users.filter(function(u){return this.indexOf(u._id)>=0;},hasNotPostedIDs);
        var hasNotPostedNames = hasNotPostedUsers.map(function(u) {return u.username;});
        return `The following members have not posted in the "${channelFromDB[0].name}" channel: ${JSON.stringify(hasNotPostedNames)}`;
    }

    async mostPosts(user, channel) {
        var channelFromDB = await this.channelsCollection.find(
            {$or: [{name: new RegExp(channel, "i"), members: user.id}, {name: new RegExp(channel, "i"), private: false}]}
            ).toArray();
        if(!channelFromDB[0]) {
            return "Sorry, I can't find that channel or you do not have access to it.";
        }

		var allMsgs = await this.messagesCollection.find({channelid: channelFromDB[0]._id}).toArray();
        var userIDs = allMsgs.map(function(m) {return m.creatorid;});
        user = userIDs.sort((a,b) =>
          userIDs.filter(v => v===a).length
        - userIDs.filter(v => v===b).length
        ).pop();
        user = await this.usersCollection.findOne({_id: user});
        return `${user.username} has posted most to the "${channelFromDB[0].name}" channel`;
    }

    async members(user, channel) {
        var channelFromDB = await this.channelsCollection.find(
            {$or: [{name: new RegExp(channel, "i"), members: user.id}, {name: new RegExp(channel, "i"), private: false}]}
            ).toArray();
        if(!channelFromDB[0]) {
            return "Sorry, I can't find that channel or you do not have access to it.";
        }

        if(channelFromDB[0].private === false) {
            var users = await this.usersCollection.find().toArray();
            var userNames = users.map(function(u) {return u.username;});
            return `Every user is a member of the "${channelFromDB[0].name}" channel, it is a public channel. Members: ${JSON.stringify(userNames)}`;
        } 
        var users = new Array();
        for (var userID of channelFromDB[0].members) {
            var user = await this.usersCollection.findOne({_id: userID});
            if(user) {
                users.push(user.username);
            }
        }
        return `The members of the "${channelFromDB[0].name}" channel are ${JSON.stringify(users)}`
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
            return `Sorry, I could not find your latest message in ${channelFromDB[0].name}`;
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
        var channelFromDB = await this.channelsCollection.find(
            {$or: [{name: new RegExp(channel, "i"), members: user.id}, {name: new RegExp(channel, "i"), private: false}]}
            ).toArray();
        if(!channelFromDB[0]) {
            return "Sorry, I can't find that channel or you do not have access to it.";
        }
        let startDate = moment(date).toDate();
        let endDate = moment(startDate).add(24, 'hours').toDate();
        let messages = await this.messagesCollection.find({channelid: channelFromDB[0]._id, createdat: {$gte : startDate, $lt: endDate}}).toArray();
        return `You posted ${messages.length} times in "${channelFromDB[0].name}" on ${date}`;
        
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