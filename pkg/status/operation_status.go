package status

type OperationStatus int

var BookmarkCreateSuccess = 100
var BookmarkCreateFailed = 101
var BookmarkUpdateSuccess = 200
var BookmarkUpdateFailed = 201
var BookmarkUpdateMatched = 202
var BookmarkUpdateUpserted = 203
