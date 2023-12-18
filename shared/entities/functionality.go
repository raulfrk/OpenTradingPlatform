package entities

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"tradingplatform/shared/logging"

	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

func HashStruct(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:]), nil
}

type Payloader interface {
	ToPayload() []byte
}

func (b *Bar) SetFingerprint() {
	b.Fingerprint, _ = HashStruct(b)
}

func (o *Orderbook) SetFingerprint() {
	o.Fingerprint, _ = HashStruct(o)
}

func (q *Quote) SetFingerprint() {
	q.Fingerprint, _ = HashStruct(q)
}

func (t *Trade) SetFingerprint() {
	t.Fingerprint, _ = HashStruct(t)
}

func (l *LULD) SetFingerprint() {
	l.Fingerprint, _ = HashStruct(l)
}

func (s *TradingStatus) SetFingerprint() {
	s.Fingerprint, _ = HashStruct(s)
}

func (n *News) SetFingerprint() {
	n.Fingerprint, _ = HashStruct(n)
}

func (b *Bar) SetSource(source string) {
	b.Source = source
}

func (o *Orderbook) SetSource(source string) {
	o.Source = source

	for i := range o.Bids {
		o.Bids[i].SetSource(source)
	}

	for i := range o.Asks {
		o.Asks[i].SetSource(source)
	}
}

func (o *OrderbookEntry) SetSource(source string) {
	o.Source = source
}

func (q *Quote) SetSource(source string) {
	q.Source = source
}

func (t *Trade) SetSource(source string) {
	t.Source = source
}

func (l *LULD) SetSource(source string) {
	l.Source = source
}

func (s *TradingStatus) SetSource(source string) {
	s.Source = source
}

func (n *News) SetSource(source string) {
	n.Source = source
}

func (b *Bar) SetExchange(exchange string) {
	b.Exchange = exchange
}

func (o *Orderbook) SetExchange(exchange string) {
	o.Exchange = exchange
}

func (q *Quote) SetExchange(exchange string) {
	q.Exchange = exchange
}

func (t *Trade) SetExchange(exchange string) {
	t.Exchange = exchange
}

type ExchangeSettable interface {
	SetExchange(string)
}

type SourceSettable interface {
	SetSource(string)
}

type ProtoReflectable interface {
	ProtoReflect() protoreflect.Message
}

type Fingerprintable interface {
	GetFingerprint() string
}

func GeneratePayload(p ProtoReflectable) []byte {
	payload, err := proto.Marshal(p)
	if err != nil {
		logging.Log().Error().Err(err).Msg("marshalling")
	}
	return payload
}

func (b *Bar) ToPayload() []byte {
	return GeneratePayload(b)
}

func (o *Orderbook) ToPayload() []byte {
	return GeneratePayload(o)
}

func (o *Quote) ToPayload() []byte {
	return GeneratePayload(o)
}

func (o *Trade) ToPayload() []byte {
	return GeneratePayload(o)
}

func (l *LULD) ToPayload() []byte {
	return GeneratePayload(l)
}

func (s *TradingStatus) ToPayload() []byte {
	return GeneratePayload(s)
}

func (n *News) ToPayload() []byte {
	return GeneratePayload(n)
}
