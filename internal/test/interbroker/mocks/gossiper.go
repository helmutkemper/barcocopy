// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	interbroker "github.com/jorgebay/soda/internal/interbroker"
	mock "github.com/stretchr/testify/mock"

	types "github.com/jorgebay/soda/internal/types"
)

// Gossiper is an autogenerated mock type for the Gossiper type
type Gossiper struct {
	mock.Mock
}

// AcceptConnections provides a mock function with given fields:
func (_m *Gossiper) AcceptConnections() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetGenerations provides a mock function with given fields: ordinal, token
func (_m *Gossiper) GetGenerations(ordinal int, token types.Token) ([]types.Generation, error) {
	ret := _m.Called(ordinal, token)

	var r0 []types.Generation
	if rf, ok := ret.Get(0).(func(int, types.Token) []types.Generation); ok {
		r0 = rf(ordinal, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Generation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, types.Token) error); ok {
		r1 = rf(ordinal, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Init provides a mock function with given fields:
func (_m *Gossiper) Init() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OpenConnections provides a mock function with given fields:
func (_m *Gossiper) OpenConnections() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterGenListener provides a mock function with given fields: listener
func (_m *Gossiper) RegisterGenListener(listener interbroker.GenListener) {
	_m.Called(listener)
}

// SendToFollowers provides a mock function with given fields: replicationInfo, topic, segmentId, body
func (_m *Gossiper) SendToFollowers(replicationInfo types.ReplicationInfo, topic types.TopicDataId, segmentId int64, body []byte) error {
	ret := _m.Called(replicationInfo, topic, segmentId, body)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.ReplicationInfo, types.TopicDataId, int64, []byte) error); ok {
		r0 = rf(replicationInfo, topic, segmentId, body)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendToLeader provides a mock function with given fields: replicationInfo, topic, body
func (_m *Gossiper) SendToLeader(replicationInfo types.ReplicationInfo, topic string, body []byte) error {
	ret := _m.Called(replicationInfo, topic, body)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.ReplicationInfo, string, []byte) error); ok {
		r0 = rf(replicationInfo, topic, body)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetGenerationAsAccepted provides a mock function with given fields: ordinal, newGen
func (_m *Gossiper) SetGenerationAsAccepted(ordinal int, newGen types.Generation) error {
	ret := _m.Called(ordinal, newGen)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, types.Generation) error); ok {
		r0 = rf(ordinal, newGen)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertGeneration provides a mock function with given fields: ordinal, existing, newGeneration
func (_m *Gossiper) UpsertGeneration(ordinal int, existing *types.Generation, newGeneration types.Generation) error {
	ret := _m.Called(ordinal, existing, newGeneration)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, *types.Generation, types.Generation) error); ok {
		r0 = rf(ordinal, existing, newGeneration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WaitForPeersUp provides a mock function with given fields:
func (_m *Gossiper) WaitForPeersUp() {
	_m.Called()
}