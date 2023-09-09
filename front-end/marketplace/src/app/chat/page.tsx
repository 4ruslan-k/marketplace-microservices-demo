"use client"

import { useState, useEffect } from 'react';
import io from 'socket.io-client';


export default function Home() {
    const [socket, setSocket] = useState<any>(null);
    useEffect(() => {
        if (!socket) {
          const updatedSocket = io('ws://localhost:4009', {
            path: '/chat/socket.io'
          });
          setSocket(updatedSocket);
    
          updatedSocket.on('connect', () => {
            console.log('Connected to Socket.IO server!');
          });


          updatedSocket.on('disconnect', () => {
            console.log('disconnect');
          });
        }
    
        return () => {
          socket?.disconnect();
        };
      }, []);

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <div className="z-10 max-w-5xl w-full items-center justify-between font-mono text-sm lg:flex">
        <h1>Chat</h1>
      </div>
     
    </main>
  )
}
