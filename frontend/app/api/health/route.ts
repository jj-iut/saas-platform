import { NextResponse } from 'next/server'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export async function GET() {
  try {
    const response = await fetch(`${API_URL}/health`, {
      cache: 'no-store',
    })

    if (!response.ok) {
      throw new Error('Backend health check failed')
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    return NextResponse.json(
      { status: 'unhealthy', error: 'Failed to connect to backend' },
      { status: 503 }
    )
  }
}
