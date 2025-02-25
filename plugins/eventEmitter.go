package plugins

import (
	"log"
	"sync"
)

type EventEmitter struct {
	// 使用 map 存储事件监听器
	listeners map[string]map[string]func(...interface{})
	mutex     sync.RWMutex // 添加互斥锁保证并发安全
}

// NewEventEmitter 创建并初始化一个新的 EventEmitter
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[string]map[string]func(...interface{})),
	}
}

// On 注册一个事件监听器
// 约定eventName为待监听插件的名字，listenerName为当前插件的名字
func (ee *EventEmitter) On(eventName string, listenerName string, listener func(...interface{})) {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()
	
	listeners, ok := ee.listeners[eventName]
	if !ok {
		// 如果没有这个事件，就创建一个
		listeners = make(map[string]func(...interface{}))
		ee.listeners[eventName] = listeners
	}
	listeners[listenerName] = listener
}

// Off 移除一个事件监听器
func (ee *EventEmitter) Off(eventName string, listenerName string) {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()
	
	if listeners, ok := ee.listeners[eventName]; ok {
		delete(listeners, listenerName)
	}
}

// Emit 触发一个事件
func (ee *EventEmitter) Emit(eventName string, args ...interface{}) {
	ee.mutex.RLock()
	defer ee.mutex.RUnlock()
	
	listeners, ok := ee.listeners[eventName]
	if ok {
		for _, listener := range listeners {
			// 使用匿名函数捕获可能的 panic
			func(l func(...interface{})) {
				defer func() {
					if r := recover(); r != nil {
						// 这里可以添加日志记录或其他处理
						// fmt.Printf("Event listener panic: %v\n", r)
						log.Printf("Event listener panic: %v\n", r)
					}
				}()
				l(args...)
			}(listener)
		}
	}
}

// HasListeners 检查特定事件是否有监听器
func (ee *EventEmitter) HasListeners(eventName string) bool {
	ee.mutex.RLock()
	defer ee.mutex.RUnlock()
	
	listeners, ok := ee.listeners[eventName]
	return ok && len(listeners) > 0
}