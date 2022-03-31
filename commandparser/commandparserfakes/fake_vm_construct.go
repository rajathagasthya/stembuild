// Code generated by counterfeiter. DO NOT EDIT.
package commandparserfakes

import (
	"sync"

	"github.com/cloudfoundry/stembuild/commandparser"
)

type FakeVmConstruct struct {
	PrepareVMStub        func() error
	prepareVMMutex       sync.RWMutex
	prepareVMArgsForCall []struct {
	}
	prepareVMReturns struct {
		result1 error
	}
	prepareVMReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeVmConstruct) PrepareVM() error {
	fake.prepareVMMutex.Lock()
	ret, specificReturn := fake.prepareVMReturnsOnCall[len(fake.prepareVMArgsForCall)]
	fake.prepareVMArgsForCall = append(fake.prepareVMArgsForCall, struct {
	}{})
	fake.recordInvocation("PrepareVM", []interface{}{})
	fake.prepareVMMutex.Unlock()
	if fake.PrepareVMStub != nil {
		return fake.PrepareVMStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.prepareVMReturns
	return fakeReturns.result1
}

func (fake *FakeVmConstruct) PrepareVMCallCount() int {
	fake.prepareVMMutex.RLock()
	defer fake.prepareVMMutex.RUnlock()
	return len(fake.prepareVMArgsForCall)
}

func (fake *FakeVmConstruct) PrepareVMCalls(stub func() error) {
	fake.prepareVMMutex.Lock()
	defer fake.prepareVMMutex.Unlock()
	fake.PrepareVMStub = stub
}

func (fake *FakeVmConstruct) PrepareVMReturns(result1 error) {
	fake.prepareVMMutex.Lock()
	defer fake.prepareVMMutex.Unlock()
	fake.PrepareVMStub = nil
	fake.prepareVMReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeVmConstruct) PrepareVMReturnsOnCall(i int, result1 error) {
	fake.prepareVMMutex.Lock()
	defer fake.prepareVMMutex.Unlock()
	fake.PrepareVMStub = nil
	if fake.prepareVMReturnsOnCall == nil {
		fake.prepareVMReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.prepareVMReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeVmConstruct) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.prepareVMMutex.RLock()
	defer fake.prepareVMMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeVmConstruct) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ commandparser.VmConstruct = new(FakeVmConstruct)
