// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/ONSdigital/dp-elasticsearch/v3/client"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"sync"
)

var (
	lockClientMockAddDocument    sync.RWMutex
	lockClientMockBulkIndexAdd   sync.RWMutex
	lockClientMockBulkIndexClose sync.RWMutex
	lockClientMockBulkUpdate     sync.RWMutex
	lockClientMockChecker        sync.RWMutex
	lockClientMockCreateIndex    sync.RWMutex
	lockClientMockDeleteIndex    sync.RWMutex
	lockClientMockDeleteIndices  sync.RWMutex
	lockClientMockGetIndices     sync.RWMutex
	lockClientMockNewBulkIndexer sync.RWMutex
	lockClientMockUpdateAliases  sync.RWMutex
)

// Ensure, that ClientMock does implement client.Client.
// If this is not the case, regenerate this file with moq.
var _ client.Client = &ClientMock{}

// ClientMock is a mock implementation of client.Client.
//
//     func TestSomethingThatUsesClient(t *testing.T) {
//
//         // make and configure a mocked client.Client
//         mockedClient := &ClientMock{
//             AddDocumentFunc: func(ctx context.Context, indexName string, documentID string, document []byte, opts *client.AddDocumentOptions) error {
// 	               panic("mock out the AddDocument method")
//             },
//             BulkIndexAddFunc: func(ctx context.Context, action client.BulkIndexerAction, index string, documentID string, document []byte) error {
// 	               panic("mock out the BulkIndexAdd method")
//             },
//             BulkIndexCloseFunc: func(in1 context.Context) error {
// 	               panic("mock out the BulkIndexClose method")
//             },
//             BulkUpdateFunc: func(ctx context.Context, indexName string, url string, settings []byte) ([]byte, error) {
// 	               panic("mock out the BulkUpdate method")
//             },
//             CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error {
// 	               panic("mock out the Checker method")
//             },
//             CreateIndexFunc: func(ctx context.Context, indexName string, indexSettings []byte) error {
// 	               panic("mock out the CreateIndex method")
//             },
//             DeleteIndexFunc: func(ctx context.Context, indexName string) error {
// 	               panic("mock out the DeleteIndex method")
//             },
//             DeleteIndicesFunc: func(ctx context.Context, indices []string) error {
// 	               panic("mock out the DeleteIndices method")
//             },
//             GetIndicesFunc: func(ctx context.Context, indexPatterns []string) ([]byte, error) {
// 	               panic("mock out the GetIndices method")
//             },
//             NewBulkIndexerFunc: func(in1 context.Context) error {
// 	               panic("mock out the NewBulkIndexer method")
//             },
//             UpdateAliasesFunc: func(ctx context.Context, alias string, removeIndices []string, addIndices []string) error {
// 	               panic("mock out the UpdateAliases method")
//             },
//         }
//
//         // use mockedClient in code that requires client.Client
//         // and then make assertions.
//
//     }
type ClientMock struct {
	// AddDocumentFunc mocks the AddDocument method.
	AddDocumentFunc func(ctx context.Context, indexName string, documentID string, document []byte, opts *client.AddDocumentOptions) error

	// BulkIndexAddFunc mocks the BulkIndexAdd method.
	BulkIndexAddFunc func(ctx context.Context, action client.BulkIndexerAction, index string, documentID string, document []byte) error

	// BulkIndexCloseFunc mocks the BulkIndexClose method.
	BulkIndexCloseFunc func(in1 context.Context) error

	// BulkUpdateFunc mocks the BulkUpdate method.
	BulkUpdateFunc func(ctx context.Context, indexName string, url string, settings []byte) ([]byte, error)

	// CheckerFunc mocks the Checker method.
	CheckerFunc func(ctx context.Context, state *healthcheck.CheckState) error

	// CreateIndexFunc mocks the CreateIndex method.
	CreateIndexFunc func(ctx context.Context, indexName string, indexSettings []byte) error

	// DeleteIndexFunc mocks the DeleteIndex method.
	DeleteIndexFunc func(ctx context.Context, indexName string) error

	// DeleteIndicesFunc mocks the DeleteIndices method.
	DeleteIndicesFunc func(ctx context.Context, indices []string) error

	// GetIndicesFunc mocks the GetIndices method.
	GetIndicesFunc func(ctx context.Context, indexPatterns []string) ([]byte, error)

	// NewBulkIndexerFunc mocks the NewBulkIndexer method.
	NewBulkIndexerFunc func(in1 context.Context) error

	// UpdateAliasesFunc mocks the UpdateAliases method.
	UpdateAliasesFunc func(ctx context.Context, alias string, removeIndices []string, addIndices []string) error

	// calls tracks calls to the methods.
	calls struct {
		// AddDocument holds details about calls to the AddDocument method.
		AddDocument []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// IndexName is the indexName argument value.
			IndexName string
			// DocumentID is the documentID argument value.
			DocumentID string
			// Document is the document argument value.
			Document []byte
			// Opts is the opts argument value.
			Opts *client.AddDocumentOptions
		}
		// BulkIndexAdd holds details about calls to the BulkIndexAdd method.
		BulkIndexAdd []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Action is the action argument value.
			Action client.BulkIndexerAction
			// Index is the index argument value.
			Index string
			// DocumentID is the documentID argument value.
			DocumentID string
			// Document is the document argument value.
			Document []byte
		}
		// BulkIndexClose holds details about calls to the BulkIndexClose method.
		BulkIndexClose []struct {
			// In1 is the in1 argument value.
			In1 context.Context
		}
		// BulkUpdate holds details about calls to the BulkUpdate method.
		BulkUpdate []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// IndexName is the indexName argument value.
			IndexName string
			// URL is the url argument value.
			URL string
			// Settings is the settings argument value.
			Settings []byte
		}
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// State is the state argument value.
			State *healthcheck.CheckState
		}
		// CreateIndex holds details about calls to the CreateIndex method.
		CreateIndex []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// IndexName is the indexName argument value.
			IndexName string
			// IndexSettings is the indexSettings argument value.
			IndexSettings []byte
		}
		// DeleteIndex holds details about calls to the DeleteIndex method.
		DeleteIndex []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// IndexName is the indexName argument value.
			IndexName string
		}
		// DeleteIndices holds details about calls to the DeleteIndices method.
		DeleteIndices []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Indices is the indices argument value.
			Indices []string
		}
		// GetIndices holds details about calls to the GetIndices method.
		GetIndices []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// IndexPatterns is the indexPatterns argument value.
			IndexPatterns []string
		}
		// NewBulkIndexer holds details about calls to the NewBulkIndexer method.
		NewBulkIndexer []struct {
			// In1 is the in1 argument value.
			In1 context.Context
		}
		// UpdateAliases holds details about calls to the UpdateAliases method.
		UpdateAliases []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Alias is the alias argument value.
			Alias string
			// RemoveIndices is the removeIndices argument value.
			RemoveIndices []string
			// AddIndices is the addIndices argument value.
			AddIndices []string
		}
	}
}

// AddDocument calls AddDocumentFunc.
func (mock *ClientMock) AddDocument(ctx context.Context, indexName string, documentID string, document []byte, opts *client.AddDocumentOptions) error {
	if mock.AddDocumentFunc == nil {
		panic("ClientMock.AddDocumentFunc: method is nil but Client.AddDocument was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		IndexName  string
		DocumentID string
		Document   []byte
		Opts       *client.AddDocumentOptions
	}{
		Ctx:        ctx,
		IndexName:  indexName,
		DocumentID: documentID,
		Document:   document,
		Opts:       opts,
	}
	lockClientMockAddDocument.Lock()
	mock.calls.AddDocument = append(mock.calls.AddDocument, callInfo)
	lockClientMockAddDocument.Unlock()
	return mock.AddDocumentFunc(ctx, indexName, documentID, document, opts)
}

// AddDocumentCalls gets all the calls that were made to AddDocument.
// Check the length with:
//     len(mockedClient.AddDocumentCalls())
func (mock *ClientMock) AddDocumentCalls() []struct {
	Ctx        context.Context
	IndexName  string
	DocumentID string
	Document   []byte
	Opts       *client.AddDocumentOptions
} {
	var calls []struct {
		Ctx        context.Context
		IndexName  string
		DocumentID string
		Document   []byte
		Opts       *client.AddDocumentOptions
	}
	lockClientMockAddDocument.RLock()
	calls = mock.calls.AddDocument
	lockClientMockAddDocument.RUnlock()
	return calls
}

// BulkIndexAdd calls BulkIndexAddFunc.
func (mock *ClientMock) BulkIndexAdd(ctx context.Context, action client.BulkIndexerAction, index string, documentID string, document []byte) error {
	if mock.BulkIndexAddFunc == nil {
		panic("ClientMock.BulkIndexAddFunc: method is nil but Client.BulkIndexAdd was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		Action     client.BulkIndexerAction
		Index      string
		DocumentID string
		Document   []byte
	}{
		Ctx:        ctx,
		Action:     action,
		Index:      index,
		DocumentID: documentID,
		Document:   document,
	}
	lockClientMockBulkIndexAdd.Lock()
	mock.calls.BulkIndexAdd = append(mock.calls.BulkIndexAdd, callInfo)
	lockClientMockBulkIndexAdd.Unlock()
	return mock.BulkIndexAddFunc(ctx, action, index, documentID, document)
}

// BulkIndexAddCalls gets all the calls that were made to BulkIndexAdd.
// Check the length with:
//     len(mockedClient.BulkIndexAddCalls())
func (mock *ClientMock) BulkIndexAddCalls() []struct {
	Ctx        context.Context
	Action     client.BulkIndexerAction
	Index      string
	DocumentID string
	Document   []byte
} {
	var calls []struct {
		Ctx        context.Context
		Action     client.BulkIndexerAction
		Index      string
		DocumentID string
		Document   []byte
	}
	lockClientMockBulkIndexAdd.RLock()
	calls = mock.calls.BulkIndexAdd
	lockClientMockBulkIndexAdd.RUnlock()
	return calls
}

// BulkIndexClose calls BulkIndexCloseFunc.
func (mock *ClientMock) BulkIndexClose(in1 context.Context) error {
	if mock.BulkIndexCloseFunc == nil {
		panic("ClientMock.BulkIndexCloseFunc: method is nil but Client.BulkIndexClose was just called")
	}
	callInfo := struct {
		In1 context.Context
	}{
		In1: in1,
	}
	lockClientMockBulkIndexClose.Lock()
	mock.calls.BulkIndexClose = append(mock.calls.BulkIndexClose, callInfo)
	lockClientMockBulkIndexClose.Unlock()
	return mock.BulkIndexCloseFunc(in1)
}

// BulkIndexCloseCalls gets all the calls that were made to BulkIndexClose.
// Check the length with:
//     len(mockedClient.BulkIndexCloseCalls())
func (mock *ClientMock) BulkIndexCloseCalls() []struct {
	In1 context.Context
} {
	var calls []struct {
		In1 context.Context
	}
	lockClientMockBulkIndexClose.RLock()
	calls = mock.calls.BulkIndexClose
	lockClientMockBulkIndexClose.RUnlock()
	return calls
}

// BulkUpdate calls BulkUpdateFunc.
func (mock *ClientMock) BulkUpdate(ctx context.Context, indexName string, url string, settings []byte) ([]byte, error) {
	if mock.BulkUpdateFunc == nil {
		panic("ClientMock.BulkUpdateFunc: method is nil but Client.BulkUpdate was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		IndexName string
		URL       string
		Settings  []byte
	}{
		Ctx:       ctx,
		IndexName: indexName,
		URL:       url,
		Settings:  settings,
	}
	lockClientMockBulkUpdate.Lock()
	mock.calls.BulkUpdate = append(mock.calls.BulkUpdate, callInfo)
	lockClientMockBulkUpdate.Unlock()
	return mock.BulkUpdateFunc(ctx, indexName, url, settings)
}

// BulkUpdateCalls gets all the calls that were made to BulkUpdate.
// Check the length with:
//     len(mockedClient.BulkUpdateCalls())
func (mock *ClientMock) BulkUpdateCalls() []struct {
	Ctx       context.Context
	IndexName string
	URL       string
	Settings  []byte
} {
	var calls []struct {
		Ctx       context.Context
		IndexName string
		URL       string
		Settings  []byte
	}
	lockClientMockBulkUpdate.RLock()
	calls = mock.calls.BulkUpdate
	lockClientMockBulkUpdate.RUnlock()
	return calls
}

// Checker calls CheckerFunc.
func (mock *ClientMock) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("ClientMock.CheckerFunc: method is nil but Client.Checker was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		State *healthcheck.CheckState
	}{
		Ctx:   ctx,
		State: state,
	}
	lockClientMockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	lockClientMockChecker.Unlock()
	return mock.CheckerFunc(ctx, state)
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//     len(mockedClient.CheckerCalls())
func (mock *ClientMock) CheckerCalls() []struct {
	Ctx   context.Context
	State *healthcheck.CheckState
} {
	var calls []struct {
		Ctx   context.Context
		State *healthcheck.CheckState
	}
	lockClientMockChecker.RLock()
	calls = mock.calls.Checker
	lockClientMockChecker.RUnlock()
	return calls
}

// CreateIndex calls CreateIndexFunc.
func (mock *ClientMock) CreateIndex(ctx context.Context, indexName string, indexSettings []byte) error {
	if mock.CreateIndexFunc == nil {
		panic("ClientMock.CreateIndexFunc: method is nil but Client.CreateIndex was just called")
	}
	callInfo := struct {
		Ctx           context.Context
		IndexName     string
		IndexSettings []byte
	}{
		Ctx:           ctx,
		IndexName:     indexName,
		IndexSettings: indexSettings,
	}
	lockClientMockCreateIndex.Lock()
	mock.calls.CreateIndex = append(mock.calls.CreateIndex, callInfo)
	lockClientMockCreateIndex.Unlock()
	return mock.CreateIndexFunc(ctx, indexName, indexSettings)
}

// CreateIndexCalls gets all the calls that were made to CreateIndex.
// Check the length with:
//     len(mockedClient.CreateIndexCalls())
func (mock *ClientMock) CreateIndexCalls() []struct {
	Ctx           context.Context
	IndexName     string
	IndexSettings []byte
} {
	var calls []struct {
		Ctx           context.Context
		IndexName     string
		IndexSettings []byte
	}
	lockClientMockCreateIndex.RLock()
	calls = mock.calls.CreateIndex
	lockClientMockCreateIndex.RUnlock()
	return calls
}

// DeleteIndex calls DeleteIndexFunc.
func (mock *ClientMock) DeleteIndex(ctx context.Context, indexName string) error {
	if mock.DeleteIndexFunc == nil {
		panic("ClientMock.DeleteIndexFunc: method is nil but Client.DeleteIndex was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		IndexName string
	}{
		Ctx:       ctx,
		IndexName: indexName,
	}
	lockClientMockDeleteIndex.Lock()
	mock.calls.DeleteIndex = append(mock.calls.DeleteIndex, callInfo)
	lockClientMockDeleteIndex.Unlock()
	return mock.DeleteIndexFunc(ctx, indexName)
}

// DeleteIndexCalls gets all the calls that were made to DeleteIndex.
// Check the length with:
//     len(mockedClient.DeleteIndexCalls())
func (mock *ClientMock) DeleteIndexCalls() []struct {
	Ctx       context.Context
	IndexName string
} {
	var calls []struct {
		Ctx       context.Context
		IndexName string
	}
	lockClientMockDeleteIndex.RLock()
	calls = mock.calls.DeleteIndex
	lockClientMockDeleteIndex.RUnlock()
	return calls
}

// DeleteIndices calls DeleteIndicesFunc.
func (mock *ClientMock) DeleteIndices(ctx context.Context, indices []string) error {
	if mock.DeleteIndicesFunc == nil {
		panic("ClientMock.DeleteIndicesFunc: method is nil but Client.DeleteIndices was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Indices []string
	}{
		Ctx:     ctx,
		Indices: indices,
	}
	lockClientMockDeleteIndices.Lock()
	mock.calls.DeleteIndices = append(mock.calls.DeleteIndices, callInfo)
	lockClientMockDeleteIndices.Unlock()
	return mock.DeleteIndicesFunc(ctx, indices)
}

// DeleteIndicesCalls gets all the calls that were made to DeleteIndices.
// Check the length with:
//     len(mockedClient.DeleteIndicesCalls())
func (mock *ClientMock) DeleteIndicesCalls() []struct {
	Ctx     context.Context
	Indices []string
} {
	var calls []struct {
		Ctx     context.Context
		Indices []string
	}
	lockClientMockDeleteIndices.RLock()
	calls = mock.calls.DeleteIndices
	lockClientMockDeleteIndices.RUnlock()
	return calls
}

// GetIndices calls GetIndicesFunc.
func (mock *ClientMock) GetIndices(ctx context.Context, indexPatterns []string) ([]byte, error) {
	if mock.GetIndicesFunc == nil {
		panic("ClientMock.GetIndicesFunc: method is nil but Client.GetIndices was just called")
	}
	callInfo := struct {
		Ctx           context.Context
		IndexPatterns []string
	}{
		Ctx:           ctx,
		IndexPatterns: indexPatterns,
	}
	lockClientMockGetIndices.Lock()
	mock.calls.GetIndices = append(mock.calls.GetIndices, callInfo)
	lockClientMockGetIndices.Unlock()
	return mock.GetIndicesFunc(ctx, indexPatterns)
}

// GetIndicesCalls gets all the calls that were made to GetIndices.
// Check the length with:
//     len(mockedClient.GetIndicesCalls())
func (mock *ClientMock) GetIndicesCalls() []struct {
	Ctx           context.Context
	IndexPatterns []string
} {
	var calls []struct {
		Ctx           context.Context
		IndexPatterns []string
	}
	lockClientMockGetIndices.RLock()
	calls = mock.calls.GetIndices
	lockClientMockGetIndices.RUnlock()
	return calls
}

// NewBulkIndexer calls NewBulkIndexerFunc.
func (mock *ClientMock) NewBulkIndexer(in1 context.Context) error {
	if mock.NewBulkIndexerFunc == nil {
		panic("ClientMock.NewBulkIndexerFunc: method is nil but Client.NewBulkIndexer was just called")
	}
	callInfo := struct {
		In1 context.Context
	}{
		In1: in1,
	}
	lockClientMockNewBulkIndexer.Lock()
	mock.calls.NewBulkIndexer = append(mock.calls.NewBulkIndexer, callInfo)
	lockClientMockNewBulkIndexer.Unlock()
	return mock.NewBulkIndexerFunc(in1)
}

// NewBulkIndexerCalls gets all the calls that were made to NewBulkIndexer.
// Check the length with:
//     len(mockedClient.NewBulkIndexerCalls())
func (mock *ClientMock) NewBulkIndexerCalls() []struct {
	In1 context.Context
} {
	var calls []struct {
		In1 context.Context
	}
	lockClientMockNewBulkIndexer.RLock()
	calls = mock.calls.NewBulkIndexer
	lockClientMockNewBulkIndexer.RUnlock()
	return calls
}

// UpdateAliases calls UpdateAliasesFunc.
func (mock *ClientMock) UpdateAliases(ctx context.Context, alias string, removeIndices []string, addIndices []string) error {
	if mock.UpdateAliasesFunc == nil {
		panic("ClientMock.UpdateAliasesFunc: method is nil but Client.UpdateAliases was just called")
	}
	callInfo := struct {
		Ctx           context.Context
		Alias         string
		RemoveIndices []string
		AddIndices    []string
	}{
		Ctx:           ctx,
		Alias:         alias,
		RemoveIndices: removeIndices,
		AddIndices:    addIndices,
	}
	lockClientMockUpdateAliases.Lock()
	mock.calls.UpdateAliases = append(mock.calls.UpdateAliases, callInfo)
	lockClientMockUpdateAliases.Unlock()
	return mock.UpdateAliasesFunc(ctx, alias, removeIndices, addIndices)
}

// UpdateAliasesCalls gets all the calls that were made to UpdateAliases.
// Check the length with:
//     len(mockedClient.UpdateAliasesCalls())
func (mock *ClientMock) UpdateAliasesCalls() []struct {
	Ctx           context.Context
	Alias         string
	RemoveIndices []string
	AddIndices    []string
} {
	var calls []struct {
		Ctx           context.Context
		Alias         string
		RemoveIndices []string
		AddIndices    []string
	}
	lockClientMockUpdateAliases.RLock()
	calls = mock.calls.UpdateAliases
	lockClientMockUpdateAliases.RUnlock()
	return calls
}
