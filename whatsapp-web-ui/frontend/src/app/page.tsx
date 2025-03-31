'use client'

import React, { useEffect, useState } from 'react'
import { Chat, Message, getChats, getMessages, sendMessage } from '@/lib/api'

export default function Home() {
  const [chats, setChats] = useState<Chat[]>([])
  const [selectedChat, setSelectedChat] = useState<Chat | null>(null)
  const [messages, setMessages] = useState<Message[]>([])
  const [newMessage, setNewMessage] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadChats()
  }, [])

  useEffect(() => {
    if (selectedChat) {
      loadMessages(selectedChat.id)
    }
  }, [selectedChat])

  async function loadChats() {
    try {
      setLoading(true)
      const fetchedChats = await getChats()
      setChats(fetchedChats)
    } catch (err) {
      setError('Failed to load chats')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  async function loadMessages(chatId: string) {
    try {
      setLoading(true)
      const fetchedMessages = await getMessages(chatId)
      setMessages(fetchedMessages)
    } catch (err) {
      setError('Failed to load messages')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  async function handleSendMessage(e: React.FormEvent) {
    e.preventDefault()
    if (!selectedChat || !newMessage.trim()) return

    try {
      const message = await sendMessage(selectedChat.id, newMessage)
      setMessages(prev => [...prev, message])
      setNewMessage('')
    } catch (err) {
      setError('Failed to send message')
      console.error(err)
    }
  }

  if (error) {
    return (
      <main className="flex min-h-screen flex-col items-center justify-center p-24">
        <div className="text-red-500 text-xl">{error}</div>
      </main>
    )
  }

  return (
    <main className="flex min-h-screen">
      {/* Chat List */}
      <div className="w-64 bg-gray-100 p-4 border-r">
        <h2 className="text-xl font-bold mb-4">Chats</h2>
        {loading ? (
          <div>Loading chats...</div>
        ) : (
          <div className="space-y-2">
            {chats.map(chat => (
              <button
                key={chat.id}
                onClick={() => setSelectedChat(chat)}
                className={`w-full text-left p-2 rounded ${
                  selectedChat?.id === chat.id ? 'bg-blue-100' : 'hover:bg-gray-200'
                }`}
              >
                {chat.name}
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Chat Messages */}
      <div className="flex-1 flex flex-col">
        {selectedChat ? (
          <>
            <div className="p-4 border-b">
              <h2 className="text-xl font-bold">{selectedChat.name}</h2>
            </div>
            <div className="flex-1 overflow-y-auto p-4 space-y-4">
              {messages.map(message => (
                <div
                  key={message.id}
                  className={`max-w-[70%] rounded-lg p-3 ${
                    message.sender === 'user'
                      ? 'bg-blue-100 ml-auto'
                      : 'bg-gray-100'
                  }`}
                >
                  {message.content}
                </div>
              ))}
            </div>
            <form onSubmit={handleSendMessage} className="p-4 border-t">
              <div className="flex space-x-2">
                <input
                  type="text"
                  value={newMessage}
                  onChange={(e) => setNewMessage(e.target.value)}
                  placeholder="Type a message..."
                  className="flex-1 p-2 border rounded"
                />
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                >
                  Send
                </button>
              </div>
            </form>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center">
            <p className="text-gray-500">Select a chat to start messaging</p>
          </div>
        )}
      </div>
    </main>
  )
} 