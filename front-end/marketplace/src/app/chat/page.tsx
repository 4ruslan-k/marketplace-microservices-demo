"use client"

import { useState, useEffect,useRef} from 'react';
import io from 'socket.io-client';
import { useFetchMessages } from '../requests/chatHooks';
import IconButton from '@mui/material/IconButton';
import Send from '@mui/icons-material/Send';


const initialMessages: Array<string> = [];


export default function Home() {
    const [socket, setSocket] = useState<any>(null);
    const [messages, setMessages] = useState<string[]>(initialMessages);
    const { refetchMessages,messages: fetchedMessages } = useFetchMessages();

    useEffect(() => {
      setMessages(fetchedMessages)
    }, [fetchedMessages])
    

    useEffect(() => {
        if (!socket) {
          const updatedSocket = io('ws://localhost:4001', {
            path: '/chat/socket.io',
            withCredentials: true,
          });
          setSocket(updatedSocket);
    
          updatedSocket.on('connect', () => {
            console.log('Connected to Socket.IO server!');
            refetchMessages()
          });


          updatedSocket.on('disconnect', () => {
            console.log('disconnect');
          });

          updatedSocket.on('newMessage', (message) => {
            console.log("newMessage -> message", message)
            setMessages((prevMessages) => [...prevMessages, message]);
          });
        }
    
        return () => {
          socket?.disconnect();
        };
      }, []);

      const sendMessage = (message: string) => {
          socket.emit('addMessage', message, (response: any) => {
            console.log("addMessage -> message", message)
            console.log("response", response)
          });
      };

  const textAreaRef = useRef<HTMLTextAreaElement>(null);
      const submit = (e: any) => {
        e.preventDefault();
        const value = textAreaRef?.current?.value;
        if (value) {
          sendMessage(value);
          textAreaRef.current.value = '';
        }
      };

      const handleKeyDown = (e: any) => {
        if (e.key === 'Enter') {
          submit(e);
        }
      };

  return (
    <main  className="min-h-screen">
      <div className="items-center justify-between p-24">
      <div className="w-full items-center justify-between text-sm">
        <div className="flex w-full items-center">
      <form className="flex appearance-none rounded-md bg-gray-800 outline-none focus:outline-none mb-24">
        <textarea
          ref={textAreaRef}
          onKeyDown={(e) => handleKeyDown(e)}
          id="minput"
          placeholder="Message"
          className="mb-2 h-64 w-320 flex-grow appearance-nonse rounded-md border-none bg-gray-800 text-white placeholder-slate-400 focus:outline-none focus:ring-transparent"
        ></textarea>
        <IconButton onClick={(e) => submit(e)} variant="plain">
          <Send color='primary' />
        </IconButton>
      </form>
    </div>
      </div>
      {messages?.map((message) => (<div key={message.id}>{message.text}</div>))}
      </div>
    </main>
  )
}
