package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/dontagr/metric/models"
)

type IHasher interface {
	setKey(key string)
	setID(id string)
	setMType(mtype string)
	setDelta(delta int64)
	setValue(value float64)
	setStringValue(value string)
	getHash() string
}

type commonHasher struct {
	key  string
	data hashData
}

type hashData struct {
	id    string
	mType string
	value string
}

func (h *commonHasher) setKey(key string) {
	h.key = key
}

func (h *commonHasher) setID(id string) {
	h.data.id = id
}

func (h *commonHasher) setMType(mtype string) {
	h.data.mType = mtype
}

func (h *commonHasher) setDelta(delta int64) {
	h.data.value = fmt.Sprintf("%d", delta)
}

func (h *commonHasher) setValue(value float64) {
	h.data.value = fmt.Sprintf("%f", value)
}

func (h *commonHasher) setStringValue(value string) {
	h.data.value = value
}

func (h *commonHasher) getHash() string {
	hmacHasher := hmac.New(sha256.New, []byte(h.key))

	hmacHasher.Write([]byte(h.data.id))
	hmacHasher.Write([]byte(h.data.mType))
	hmacHasher.Write([]byte(h.data.value))

	return hex.EncodeToString(hmacHasher.Sum(nil))
}

type Manager struct {
	hasher IHasher
}

func NewHashManager() *Manager {
	return &Manager{hasher: &commonHasher{}}
}

func (h *Manager) SetKey(key string) {
	h.hasher.setKey(key)
}

func (h *Manager) SetMetrics(metric *models.Metrics) {
	h.hasher.setID(metric.ID)
	h.hasher.setMType(metric.MType)
	if metric.Delta != nil {
		h.hasher.setDelta(*metric.Delta)
	}
	if metric.Value != nil {
		h.hasher.setValue(*metric.Value)
	}
}

func (h *Manager) SetStringValue(value string) {
	h.hasher.setStringValue(value)
}

func (h *Manager) GetHash() string {
	return h.hasher.getHash()
}
