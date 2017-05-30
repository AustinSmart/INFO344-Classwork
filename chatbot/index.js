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

function handlePreReqs(req, res, witaiData) {
	let course = witaiData.entities.course[0].value;
	let coursePrereqs = prereqs[course];
	res.send(`The pre-requisites for ${course} are ${coursePrereqs.join(', ')}.`); 
}

app.use(bodyParser.text());
app.post("/v1/bot", (req, res, next) => {
	let q = req.body
	witaiClient.message(q)
	.then(data => {
		res.send(data);
		// switch (data.entities.intent[0].value) {
		// 	default:
		// 		res.send(data);
		// }
	})	
	.catch(next);
});

app.listen(port, host, () => {
	console.log(`server is listening at http://${host}:${port}`);
});
