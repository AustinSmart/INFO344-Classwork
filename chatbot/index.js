"use strict";

const express = require('express');
const { Wit } = require('node-wit');
const bodyParser = require('body-parser')

const app = express();
const port = process.env.PORT || '80';
const host = process.env.HOST || '';
const witaiToken = process.env.WITAITOKEN;
if (!witaiToken) {
	console.error("please set WITAITOKEN to your wit.ai app token");
	process.exit(1);
}

const witaiClient = new Wit({ accessToken: witaiToken });

function handleIntents(req, res, witaiData) {
	let user = req.headers["User"];
	let intents = witaiData.entities.intent;
	console.log(intents[0]);
	switch (intents[0].value) {
		//Who has made the most posts to the XYZ channel, Who is in the XYZ channel?
		case "Who":
			res.send("who made the most posts, or who is in a channel");
			break;
		//Who hasn't posted to the XYZ channel?
		case "hasn't":
		case "who hasnt":
			res.send("who hasnt posted to a channel");
			break;
		//When was my last post?, When was my last post to the XYZ channel?
		case "When":
			res.send("when was my last post, ever or in a channel");
			break;
		//How many posts have I made to the XYZ channel?, How many posts did I make to the XYZ channel yesterday?
		case "many":
		case "How many":
			res.send("how mnay posts have I made to a channel, total or in date");
			break;
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

app.listen(port, host, () => {
	console.log(`server is listening at http://${host}:${port}`);
});
