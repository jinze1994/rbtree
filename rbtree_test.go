package rbtree

import (
	"math/rand"
	"sort"
	"testing"
)

// ============================= 测试红黑树的正确性 ============================

const count = 1000

type Int int

func (i Int) Less(than Item) bool {
	return int(i) < int(than.(Int))
}

func TestCorrect(t *testing.T) {
	assert := func(ok bool, args ...interface{}) {
		if !ok {
			t.Fatal(args...)
		}
	}

	// 新建一棵红黑树
	rbt := NewRbtree()
	assert(rbt != nil)

	// 创建一个长度为 count 的整数数组, 取值范围 [0, 0.7*count)
	intArray := make([]int, count, count)
	for i := 0; i < count; i++ {
		intArray[i] = rand.Intn(int(float64(count) * 0.7))
	}
	// 用 hasInsert map 来同步记录哪些元素已被插入
	hasInsert := make(map[int]int)

	for i := 0; i < count; i++ {
		// 同步测试插入
		_, mapok := hasInsert[intArray[i]]
		hasInsert[intArray[i]] = 1
		node, ok := rbt.Insert(Int(intArray[i]))

		// 验证插入重复元素反应
		assert(mapok != ok)
		// 验证返回节点值是否正确
		assert(intArray[i] == int(node.Item.(Int)))
		// 验证Count函数
		assert(len(hasInsert) == rbt.Count())

		// 验证红黑树结构
		rbt.testStructure()
	}
	t.Log("红黑树中已正确插入", rbt.Count(), "个元素")

	sortedIntArray := make([]int, count, count)
	copy(sortedIntArray, intArray)
	sort.Ints(sortedIntArray)
	// 对 sortedIntArray 去重
	j := 0
	for i := 0; i < len(sortedIntArray); i++ {
		if sortedIntArray[j] != sortedIntArray[i] {
			j++
			sortedIntArray[j] = sortedIntArray[i]
		}
	}
	sortedIntArray = sortedIntArray[:j+1]

	assert(rbt.Count() == len(sortedIntArray), len(sortedIntArray))
	node := rbt.First()
	for i := 0; i < rbt.Count(); i++ {
		assert(int(node.Item.(Int)) == sortedIntArray[i])
		node = node.Next()
	}
	t.Log("红黑树内容 check 成功")

	for i := 0; i < count; i++ {
		_, mapok := hasInsert[i]
		node := rbt.Get(Int(i))

		// 测试Get函数反应
		if mapok {
			assert(int(node.Item.(Int)) == int(i))
		} else {
			assert(node == nil)
		}
		rbt.testStructure()
	}
	t.Log("红黑数对所有元素Get测试成功")

	for i := 0; i < count; i++ {
		// 同步测试删除
		_, mapok := hasInsert[intArray[i]]
		delete(hasInsert, intArray[i])
		item, ok := rbt.Remove(Int(intArray[i]))

		// 验证删除可重复元素的正确性
		assert(mapok == ok)
		// 验证Count函数
		assert(len(hasInsert) == rbt.Count())
		// 验证删除后节点的返回值
		if ok {
			assert(int(item.(Int)) == intArray[i])
		} else {
			assert(item == nil)
		}

		// 验证红黑树结构正确性
		rbt.testStructure()
	}
	t.Log("红黑树删除测试成功")
}

// ============================== 对比测试 Map 的速度 ==========================

var m map[int]bool

const N = 10000000
const limit = 1 << 22

// 测试 map 的插入速度
func BenchmarkMapInsert(b *testing.B) {
	m = make(map[int]bool)
	for i := 0; i < N; i++ {
		m[rand.Intn(limit)] = true
	}
	// 输出被构建的 map 的大小
	b.Log("被插入 map 中的元素", len(m))
}

// 测试 map 的查询速度
func BenchmarkMapFind(b *testing.B) {
	count := 0
	for i := 0; i < N; i++ {
		_, ok := m[rand.Intn(limit)]
		if ok == true {
			count++
		}
	}
	// 输出 map 中被找到的元素的个数
	b.Log("在 map 中被找到的个数", count)
}

// 测试 map 的删除速度
func BenchmarkMapDelete(b *testing.B) {
	count := 0
	for i := 0; i < N; i++ {
		n := rand.Intn(limit)

		_, ok := m[n]
		if ok == true {
			count++
		}

		delete(m, n)
	}
	// 输出 map 中被删除的元素的个数
	b.Log("在 map 中被删除的个数", count)
}

var rbt *Rbtree

func BenchmarkRbtInsert(b *testing.B) {
	rbt = NewRbtree()
	count := 0
	for i := 0; i < N; i++ {
		_, ok := rbt.Insert(Int(rand.Intn(limit)))
		if ok {
			count++
		}
	}
	b.Log("在红黑树中被插入的个数", count)
}

func BenchmarkRbtFind(b *testing.B) {
	count := 0
	for i := 0; i < N; i++ {
		ret := rbt.Get(Int(rand.Intn(limit)))
		if ret != nil {
			count++
		}
	}
	b.Log("在红黑树中被找到的个数", count)
}

func BenchmarkRbtRemove(b *testing.B) {
	count := 0
	for i := 0; i < N; i++ {
		_, ok := rbt.Remove(Int(rand.Intn(limit)))
		if ok {
			count++
		}
	}
	b.Log("在红黑树中被删除的个数", count)
}
