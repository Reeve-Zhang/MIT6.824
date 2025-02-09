package kvraft

import (
	"crypto/rand"
	"math/big"
	"time"

	"6.824/labrpc"
)

type Clerk struct {
	servers []*labrpc.ClientEnd
	// You will have to modify this struct.
	clientID int64
	seqNum   int
	leaderID int
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(servers []*labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.servers = servers
	// You'll have to add code here.
	ck.clientID = nrand()
	ck.seqNum = 1
	ck.leaderID = -1
	return ck
}

//
// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
func (ck *Clerk) Get(key string) string {

	// You will have to modify this function.
	var ok bool
	var leaderID int
	var lastReplyErrLeader bool
	getArgs := GetArgs{key, ck.clientID, ck.seqNum}
	for {
		getReply := GetReply{}
		leaderID = ck.leaderID
		if lastReplyErrLeader || leaderID == -1 {
			leaderID = ck.getRandServer()
		}
		ok = ck.servers[leaderID].Call("KVServer.Get", &getArgs, &getReply)
		if !ok || getReply.Err == "ErrWrongLeader" {
			lastReplyErrLeader = true
			continue
		}

		PrettyDebug(dClient, "Client %d Send GET to Server %d, GetArgs%s, and receive GetReply%s ", ck.clientID, leaderID, getArgs.String(), getReply.String())

		switch getReply.Err {
		case "OK":
			ck.seqNum++
			ck.leaderID = leaderID
			lastReplyErrLeader = false
			PrettyDebug(dClient, "Client %d Send GET to Server %d successfully", ck.clientID, leaderID)
			return getReply.Value
		case "ErrNoKey":
			ck.seqNum++
			ck.leaderID = leaderID
			lastReplyErrLeader = false
			return ""
		}
		time.Sleep(10 * time.Millisecond)
	}
	return ""
}

//
// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.PutAppend", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
func (ck *Clerk) PutAppend(key string, value string, op string) {

	// You will have to modify this function.
	var ok bool
	var leaderID int
	var lastReplyErrLeader bool
	putAppendArgs := PutAppendArgs{key, value, op, ck.clientID, ck.seqNum}
	for {
		putAppendReply := PutAppendReply{}
		leaderID = ck.leaderID
		if lastReplyErrLeader || leaderID == -1 {
			leaderID = ck.getRandServer()
		}
		ok = ck.servers[leaderID].Call("KVServer.PutAppend", &putAppendArgs, &putAppendReply)
		if !ok || putAppendReply.Err == "ErrWrongLeader" {
			lastReplyErrLeader = true
			continue
		}
		PrettyDebug(dClient, "Client %d Send PUTAPPEND to Server %d, send args %s and receive %s ", ck.clientID, leaderID, putAppendArgs.String(), putAppendReply.String())

		switch putAppendReply.Err {
		case "OK":
			PrettyDebug(dClient, "Client %d Send PUTAPPEND to Server %d successfully", ck.clientID, leaderID)
			lastReplyErrLeader = false
			ck.leaderID = leaderID
			ck.seqNum++
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
}

func (ck *Clerk) getRandServer() int {
	newServer := int(nrand()) % len(ck.servers)
	return newServer
}
