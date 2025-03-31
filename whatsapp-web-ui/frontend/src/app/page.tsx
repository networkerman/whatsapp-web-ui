'use client'

import React from 'react'
import { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'WhatsApp Web Interface',
  description: 'A web interface for WhatsApp MCP',
}

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <div className="z-10 max-w-5xl w-full items-center justify-between font-mono text-sm">
        <h1 className="text-4xl font-bold mb-8">WhatsApp Web Interface</h1>
        <p className="text-lg mb-4">Welcome to your WhatsApp Web Interface powered by Claude API.</p>
      </div>
    </main>
  )
} 