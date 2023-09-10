"use client"

import { useState, useEffect,useRef} from 'react';
import io from 'socket.io-client';
import { useFetchMessages } from '../requests/chatHooks';


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
      <div className="flex flex-col items-center justify-between p-24">
      <div className="z-10 max-w-5xl w-full items-center justify-between font-mono text-sm lg:flex">
        <div className="flex w-full items-center">
      <form className="flex appearance-none rounded-md bg-gray-800 outline-none focus:outline-none">
        <textarea
          ref={textAreaRef}
          onKeyDown={(e) => handleKeyDown(e)}
          id="minput"
          placeholder="Message"
          className="mb-2 max-h-16 flex-grow appearance-nonse rounded-md border-none bg-gray-800 text-white placeholder-slate-400 focus:outline-none focus:ring-transparent"
        ></textarea>
        <button
          onClick={(e) => submit(e)}
          className="self-end p-2 text-slate-400"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="h-4 w-4 bg-gray-800"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
            />
          </svg>
        </button>
      </form>
    </div>
      </div>
      </div>
      {messages?.map((message) => (<div key={message.id}>{message.text}</div>))}
    </main>
  )
}
