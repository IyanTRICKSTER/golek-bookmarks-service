package status

type OperationStatus int

const (
	BookmarkCreateSuccess      OperationStatus = 100
	BookmarkCreateFailed       OperationStatus = 101
	BookmarkUpdateSuccess      OperationStatus = 200
	BookmarkUpdateFailed       OperationStatus = 201
	BookmarkUpdateMatched      OperationStatus = 202
	BookmarkUpdateUpserted     OperationStatus = 203
	BookmarkDeleteFailed       OperationStatus = 300
	BookmarkDeleteSuccess      OperationStatus = 301
	OperationUnauthorized      OperationStatus = 500
	OperationAuthorized        OperationStatus = 501
	OperationForbidden         OperationStatus = 502
	BookmarkFetchingFailed     OperationStatus = 600
	BookmarkNotExist           OperationStatus = 601
	BookmarkPostFailed         OperationStatus = 602
	BookmarkPostSuccess        OperationStatus = 603
	BookmarkMatchedNotModified OperationStatus = 604
	OperationSuccess           OperationStatus = 700
	BookmarkDeletePostFailed   OperationStatus = 800
	BookmarkDeletePostSuccess  OperationStatus = 801
	BookmarkDuplicationOccurs  OperationStatus = 802
	BookmarkPostRevokeFailed   OperationStatus = 803
)

func Is(status OperationStatus, target OperationStatus) bool {
	if status == target {
		return true
	}
	return false
}
