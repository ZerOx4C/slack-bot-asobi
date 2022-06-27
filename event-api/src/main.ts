function doPost(e: GoogleAppsScript.Events.DoPost): GoogleAppsScript.Content.TextOutput {
	const json = JSON.parse(e.postData.contents);

	if (json.token != VERIFICATION_TOKEN) {
		throw new Error("invalid token.");
	}

	if (json.type == "url_verification") {
		return ContentService.createTextOutput(json.challenge);
	}

	return ContentService.createTextOutput("yay!");
}
