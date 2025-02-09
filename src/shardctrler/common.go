package shardctrler

import "fmt"

//
// Shard controler: assigns shards to replication groups.
//
// RPC interface:
// Join(servers) -- add a set of groups (gid -> server-list mapping).
// Leave(gids) -- delete a set of groups.
// Move(shard, gid) -- hand off one shard from current owner to gid.
// Query(num) -> fetch Config # num, or latest config if num==-1.
//
// A Config (configuration) describes a set of replica groups, and the
// replica group responsible for each shard. Configs are numbered. Config
// #0 is the initial configuration, with no groups and all shards
// assigned to group 0 (the invalid group).
//
// You will need to add fields to the RPC argument structs.
//

// The number of shards.
const NShards = 10

// A configuration -- an assignment of shards to groups.
// Please don't change this.
type Config struct {
	Num    int              // config number
	Shards [NShards]int     // shard -> gid
	Groups map[int][]string // gid -> servers[]
}

func (config *Config) String() string {
	return fmt.Sprintf("Config: Num: %v, Shards: %v, Groups: %v", config.Num, config.Shards, config.Groups)
}

const (
	OK             = "OK"
	ErrWrongLeader = "ErrWrongLeader"
)

type Err string

type JoinArgs struct {
	Servers map[int][]string // new GID -> servers mappings
}

type JoinReply struct {
	WrongLeader bool
	Err         Err
}

type LeaveArgs struct {
	GIDs []int
}

type LeaveReply struct {
	WrongLeader bool
	Err         Err
}

type MoveArgs struct {
	Shard int
	GID   int
}

type MoveReply struct {
	WrongLeader bool
	Err         Err
}

type QueryArgs struct {
	Num int // desired config number
}

type QueryReply struct {
	WrongLeader bool
	Err         Err
	Config      Config
}

type CommandArgs struct {
	ArgsType int              // JoinArgs, LeaveArgs, MoveArgs, QueryArgs
	Servers  map[int][]string // JoinArgs		// new GID -> servers mappings
	GIDs     []int            // LeaveArgs
	Shard    int              // MoveArgs
	GID      int              // MoveArgs
	Num      int              // QueryArgs 		// desired config number
	SeqNum   int              // for deduplication
	ClientID int64            // for deduplication
}

func (args *CommandArgs) String() string {
	return fmt.Sprintf("CommandArgs: ArgsType %v, Servers: %v, GIDs: %v, Shard: %v, GID: %v, Num: %v, SeqNum: %v, ClientID: %v", args.ArgsType, args.Servers, args.GIDs, args.Shard, args.GID, args.Num, args.SeqNum, args.ClientID)
}

type CommandReply struct {
	ReplyType   int    // JoinArgs, LeaveArgs, MoveArgs, QueryArgs
	WrongLeader bool   // Join, Leave, Move, Query
	Err         Err    // Join, Leave, Move, Query
	Config      Config // Query
}

func (reply *CommandReply) String() string {
	return fmt.Sprintf("CommandReply: ReplyType %v, WrongLeader: %v, Err: %v, Config: %v", reply.ReplyType, reply.WrongLeader, reply.Err, reply.Config)
}

const (
	JoinArgsType  = 1
	LeaveArgsType = 2
	MoveArgsType  = 3
	QueryArgsType = 4
)
