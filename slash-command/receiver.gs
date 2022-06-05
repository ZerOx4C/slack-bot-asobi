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
  const json = JSON.stringify(Object.assign({uuid: uuid}, param));
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
