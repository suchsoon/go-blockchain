package state

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/volodymyrprokopyuk/go-blockchain/chain"
	"github.com/volodymyrprokopyuk/go-blockchain/chain/account"
	"github.com/volodymyrprokopyuk/go-blockchain/chain/store"
)

type State struct {
  balances map[chain.Address]uint64
  nonces map[chain.Address]uint64
  lastBlock store.Block
  txs map[chain.Hash]chain.SigTx
  Pending *State
}

func (s *State) Nonce(addr chain.Address) uint64 {
  return s.nonces[addr]
}

func NewState(gen store.Genesis) *State {
  return &State{
    balances: maps.Clone(gen.Balances),
    nonces: make(map[chain.Address]uint64),
    txs: make(map[chain.Hash]chain.SigTx),
    Pending: &State{
      balances: maps.Clone(gen.Balances),
      nonces: make(map[chain.Address]uint64),
      txs: make(map[chain.Hash]chain.SigTx),
    },
  }
}

func (s *State) Clone() *State {
  return &State{
    balances: maps.Clone(s.balances),
    nonces: maps.Clone(s.nonces),
    lastBlock: s.lastBlock,
    txs: maps.Clone(s.txs),
    Pending: &State{
      txs: maps.Clone(s.Pending.txs),
    },
  }
}

func (s *State) Apply(sta *State) {
  s.balances = sta.balances
  s.nonces = sta.nonces
  s.lastBlock = sta.lastBlock
}

func (s *State) ResetPending() {
  s.Pending = &State{
    balances: maps.Clone(s.balances),
    nonces: maps.Clone(s.nonces),
    txs: make(map[chain.Hash]chain.SigTx),
  }
}

func (s *State) String() string {
  var bld strings.Builder
  bld.WriteString("Balances\n")
  for acc, bal := range s.balances {
    bld.WriteString(fmt.Sprintf("  %.7s: %29d\n", acc, bal))
  }
  bld.WriteString("Nonces\n")
  for acc, nonce := range s.nonces {
    bld.WriteString(fmt.Sprintf("  %.7s: %35d\n", acc, nonce))
  }
  bld.WriteString("Last block\n")
  bld.WriteString(fmt.Sprintf("  %v", s.lastBlock))
  if s.Pending != nil && len(s.Pending.txs) > 0 {
    bld.WriteString("Pending txs\n")
    for _, tx := range s.Pending.txs {
      bld.WriteString(fmt.Sprintf("  %v\n", tx))
    }
  }
  if s.Pending != nil && len(s.Pending.balances) > 0 {
    bld.WriteString("Pending balances\n")
    for acc, bal := range s.Pending.balances {
      bld.WriteString(fmt.Sprintf("  %.7s: %29d\n", acc, bal))
    }
  }
  if s.Pending != nil && len(s.Pending.nonces) > 0 {
    bld.WriteString("Pending nonces\n")
    for acc, nonce := range s.Pending.nonces {
      bld.WriteString(fmt.Sprintf("  %.7s: %35d\n", acc, nonce))
    }
  }
  return bld.String()
}

func (s *State) ApplyTx(stx chain.SigTx) error {
  hash := stx.Hash()
  valid, err := account.Verify(stx)
  if err != nil {
    return err
  }
  if !valid {
    return fmt.Errorf("%.7s: invalid signature", hash)
  }
  if stx.Nonce != s.nonces[stx.From] + 1 {
    return fmt.Errorf("%.7s: invalid nonce", hash)
  }
  if s.balances[stx.From] < stx.Value {
    return fmt.Errorf("%.7s: insufficient funds", hash)
  }
  s.balances[stx.From] -= stx.Value
  s.balances[stx.To] += stx.Value
  s.nonces[stx.From]++
  s.txs[hash] = stx
  return nil
}

func (s *State) CreateBlock() store.Block {
  pndTxs := make([]chain.SigTx, 0, len(s.Pending.txs))
  for _, tx := range s.Pending.txs {
    pndTxs = append(pndTxs, tx)
  }
  slices.SortFunc(pndTxs, func(a, b chain.SigTx) int {
    cmp := strings.Compare(string(a.From), string(b.From))
    if cmp != 0 {
      return cmp
    }
    return int(a.Nonce) - int(b.Nonce)
  })
  txs := make([]chain.SigTx, 0, len(pndTxs))
  for _, tx := range pndTxs {
    err := s.ApplyTx(tx)
    if err != nil {
      fmt.Printf("REJECTED %v\n", err)
      continue
    }
    txs = append(txs, tx)
  }
  blk := store.NewBlock(s.lastBlock.Number + 1, s.lastBlock.Hash(), txs)
  return blk
}

func (s *State) ApplyBlock(blk store.Block) error {
  if blk.Number != s.lastBlock.Number + 1 {
    return fmt.Errorf("%.7s: invalid block number", blk.Hash())
  }
  if blk.Parent != s.lastBlock.Hash() {
    return fmt.Errorf("%.7s: invalid parent hash", blk.Hash())
  }
  for _, tx := range blk.Txs {
    err := s.ApplyTx(tx)
    if err != nil {
      return err
    }
  }
  s.lastBlock = blk
  return nil
}

func (s *State) ReadBlocks(blockStoreDir string) error {
  blocks, closeBlocks, err := store.ReadBlocks(blockStoreDir)
  if err != nil {
    return err
  }
  defer closeBlocks()
  for err, blk := range blocks {
    if err != nil {
      return err
    }
    clo := s.Clone()
    err = clo.ApplyBlock(blk)
    if err != nil {
      return err
    }
    s.Apply(clo)
    s.ResetPending()
  }
  return nil
}
