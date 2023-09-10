import io from 'socket.io-client';
import React, { createContext, useState, useEffect } from 'react';

interface SocketContextProps {
  socket: any | null;
  setUserStatus?: React.Dispatch<React.SetStateAction<boolean>>;
}

export const SocketContext = createContext<SocketContextProps>({ socket: null });
export function SocketProvider({ children }: any): any {
  const [socket, setSocket] = useState(null);
  const [isUser, setUserStatus] = useState(false);

  useEffect(() => {
    if (!socket && isUser) {
      const updatedSocket = io('ws://localhost:4001', {
        jsonp: false,
        forceNew: false,
      });
      setSocket(updatedSocket);

      updatedSocket.on('connect', () => {
        console.log('Connected to Socket.IO server!');
      });
      updatedSocket.on(
        'disconnect',
        () => {
          console.log('disconnect');
        },
        []
      );
    }

    if (!isUser && socket) {
      socket?.disconnect();
      setSocket(null);
    }

    return () => {
      socket?.disconnect();
    };
  }, [isUser, socket]);

  return (
    <SocketContext.Provider
      // eslint-disable-next-line react/jsx-no-constructed-context-values
      value={{ socket, setUserStatus }}
    >
      {children}
    </SocketContext.Provider>
  );
}
