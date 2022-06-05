const LOG_CAPACITY = 20;

var sheetTable = {};

function doGet(e) {
  const response = ContentService.createTextOutput();
  response.setMimeType(ContentService.MimeType.JSON);
  response.setContent(JSON.stringify(getQueue()));
  return response;
}

function doPost(e) {
  const response = ContentService.createTextOutput();

  if (e.parameter.uuid) {
    doneQueue(e.parameter.uuid);

    response.setMimeType(ContentService.MimeType.TEXT);
    response.setContent("ok");
  }
  else {
    pushQueue(e.parameter);

    response.setMimeType(ContentService.MimeType.JSON);
    response.setContent(JSON.stringify({
      text: "",
      response_type: "in_channel",
    }));
  }

  return response;
}

function getQueue() {
  const queueSheet = obtainSheet("queue");
  if (!queueSheet) {
    return [];
  }

  const queueLength = queueSheet.getLastRow();
  if (queueLength == 0) {
    return [];
  }

  return queueSheet.getRange(1, 2, queueLength).getValues().map(values => JSON.parse(values[0]));
}

function pushQueue(param) {
  const queueSheet = obtainSheet("queue");
  if (!queueSheet) {
    return;
  }

  const uuid = Utilities.getUuid();
  const json = JSON.stringify({
    uuid       : uuid,
    channelId  : param.channel_id,
    userId     : param.user_id,
    command    : param.command,
    text       : param.text,
    responseUrl: param.response_url,
  });

  queueSheet.appendRow([uuid, json]);
}

function doneQueue(id) {
  const queueSheet = obtainSheet("queue");
  if (!queueSheet) {
    return;
  }

  const rowIndex = queueSheet.getRange(1, 1, queueSheet.getLastRow()).getValues().findIndex(values => values[0] == id);
  if (rowIndex < 0) {
    return;
  }

  queueSheet.deleteRow(rowIndex + 1);
}

function log(...args) {
  const logSheet = obtainSheet("log");
  if (!logSheet) {
    return;
  }

  logSheet.appendRow([Date.now(), ...args]);

  const logLength = logSheet.getLastRow();
  if (logLength < LOG_CAPACITY) {
    return;
  }

  logSheet.deleteRows(1, logLength - LOG_CAPACITY);
}

function obtainSheet(sheetName) {
  let sheet = sheetTable[sheetName];
  if (sheet) {
    return sheet;
  }

  const spreadsheet = SpreadsheetApp.getActiveSpreadsheet();
  sheet = spreadsheet.getSheetByName(sheetName);
  if (sheet) {
    sheetTable[sheetName] = sheet;
    return sheet;
  }

  sheet = spreadsheet.insertSheet(sheetName);
  if (sheet) {
    sheetTable[sheetName] = sheet;
    return sheet;
  }

  return null;
}

function test() {
  done("6c5d3410-d330-4e39-b036-d14199b695ae");
}
