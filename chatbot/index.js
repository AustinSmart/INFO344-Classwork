"use strict";

const express = require('express');
const { Wit } = require('node-wit');
const bodyParser = require('body-parser')
const MongoClient = require('mongodb').MongoClient;

const ChatInterface = require('./models/chat/mongointerface.js');

const app = express();
const port = process.env.PORT || '80';
const host = process.env.HOST || '';
const witaiToken = process.env.WITAITOKEN;

var chat;

if (!witaiToken) {
	console.error("please set WITAITOKEN to your wit.ai app token");
	process.exit(1);
}

const witaiClient = new Wit({ accessToken: witaiToken });

async function handleIntents(req, res, witaiData) {
	let user = JSON.parse(req.headers["user"]);
	let entities = witaiData.entities;
	console.log(JSON.stringify(entities));
	let intents = entities.intent;
	let message = entities.hasOwnProperty("message");
	let channel = entities.hasOwnProperty("channel");
	let most = entities.hasOwnProperty("most");
	let date = entities.hasOwnProperty("datetime");
	const defaultResponse = "Sorry, I'm not sure what you are asking. Try asking with different phrasing";
	if(intents[0]) {
		switch (intents[0].value) {
			//Who hasn't posted to the XYZ channel?
			case "hasn't":
			case "who hasnt":
				if (message && channel) {
					res.send(chat.hasNotPosted(entities.channel[0].value));
				} else {
					res.send(defaultResponse + " Maybe try specifying as channel");
				}
				break;
			//Who has made the most posts to the XYZ channel, Who is in the XYZ channel?
			case "Who":
				if(message && most && channel) {
					res.send(chat.mostPosts(entities.channel[0].value))
				} else if(channel) {
					res.send(chat.members(centities.channel[0].value));
				} else {
					res.send(defaultResponse + " Maybe try specifying a channel");
				}
				break;
			//When was my last post?, When was my last post to the XYZ channel?
			case "When":
				if(message && channel) {
					res.send(await chat.lastPostInChannel(user, entities.channel[0].value));
				} else if (message) {
					res.send(await chat.lastPost(user));
				} else {
					res.send(defaultResponse);
				}
				break;
			//How many posts have I made to the XYZ channel?, How many posts did I make to the XYZ channel yesterday?
			case "many":
			case "How many":
				if(message && channel && date) {
					res.send(await chat.totalMessagesInChannelOnDate(user, entities.channel[0].value, entities.datetime.value))
				} else if(message && channel){
					res.send(await chat.totalMessagesInChannel(user, entities.channel[0].value));
				} else if(message) {
					res.send(await chat.totalMessages(user));
				} else {
					res.send(defaultResponse);
				}
				break;
			default:
				res.Send(defaultResponse);
		} 
	}
}

app.use(bodyParser.text());

app.post("/v1/bot", (req, res, next) => {
	let q = req.body
	witaiClient.message(q)
	.then(data => {
		handleIntents(req, res, data);
	})	
	.catch(next);
});

var messagesUrl = 'mongodb://mongo/messages-db';
MongoClient.connect(messagesUrl).then(db =>  {
  console.log("Connected successfully to messages-db");
  let messagesCollection = db.collection('messages');
  let channelsCollection = db.collection('channels');
  chat = new ChatInterface(messagesCollection, channelsCollection);
  app.listen(port, host, () => {
		console.log(`server is listening at http://${host}:${port}`);
	});
}).catch(err => {
	console.log(err);
});