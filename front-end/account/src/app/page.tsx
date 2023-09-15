"use client";
import { useFetchUser } from './requests/userHooks';
import config from './config';
import { useRouter } from 'next/navigation';


export default function Home() {
  const { user, isLoadingUser } = useFetchUser();
  const router = useRouter();
  
  if (!isLoadingUser && user) window.location = config.MARKETPLACE_APP_URL
  if (!isLoadingUser && !user) router.push('/signin')
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
     <h1>Accounts</h1>
    </main>
  )
}
