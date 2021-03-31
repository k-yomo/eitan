import React from 'react'
import Link from 'next/link'
import { useCurrentAccountQuery } from "../generated/graphql"
import { routes } from "@src/constants/routes"

export default function Header() {
  return (
    <header>
      <div className="relative bg-white">
        <div className="mx-auto px-4 sm:px-6">
          <div
            className="flex justify-between items-center border-b-2 border-gray-100 py-4 md:justify-start md:space-x-10">
            <div className="flex justify-start lg:w-0 lg:flex-1">
              <Link href="/">
                <a>
                  <span className="sr-only">Home</span>
                  <img
                    className="h-8 w-auto sm:h-10"
                    src="/logo.png"
                    alt="サイトロゴ"
                  />
                </a>
              </Link>
            </div>
            <UserInfo/>
          </div>
        </div>
      </div>
    </header>
  )
}

const UserInfo = () => {
  const { data, loading, error } = useCurrentAccountQuery()
  if (loading) {
    return <></>
  }

  if (error) {
    console.log(error)
  }

  return (
    <>
      {data && (data.currentAccount.screenImgUrl ? (
        <div className="flex -space-x-1 overflow-hidden">
          <img className="inline-block h-10 w-10 rounded-full ring-2 ring-white" src={data.currentAccount.screenImgUrl}/>
        </div>
      ) : (
        <div>
          {data.currentAccount.displayName}
        </div>
      ))}
      {(!data || !data.currentAccount) && (
        <>
          <div className="md:flex items-center justify-end md:flex-1 lg:w-0">
            <Link href={routes.login()}>
              <a
                className="whitespace-nowrap text-base font-medium text-gray-500 hover:text-gray-900 transition duration-300">
                Log in
              </a>
            </Link>
            <Link href={routes.signUp()}>
              <a
                className="ml-8 whitespace-nowrap inline-flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-gray-800 hover:bg-white hover:border-gray-800 hover:text-gray-900 transition duration-300">
                Sign Up
              </a>
            </Link>
          </div>
        </>
      )}
    </>
  )
}