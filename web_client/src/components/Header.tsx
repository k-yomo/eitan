import React from 'react';
import Link from 'next/link';

export default function Header() {
  return (
    <header>
      <div className="relative bg-white">
        <div className="mx-auto px-4 sm:px-6">
          <div className="flex justify-between items-center border-b-2 border-gray-100 py-4 md:justify-start md:space-x-10">
            <div className="flex justify-start lg:w-0 lg:flex-1">
              <Link href="/">
                <a>
                  <span className="sr-only">ホーム</span>
                  <img
                    className="h-8 w-auto sm:h-10"
                    src="/logo.png"
                    alt="サイトロゴ"
                  />
                </a>
              </Link>
            </div>
            <div className="md:flex items-center justify-end md:flex-1 lg:w-0">
              <Link
                href="/login"

              >
                <a className="whitespace-nowrap text-base font-medium text-gray-500 hover:text-gray-900 transition duration-300">
                  ログイン
                </a>
              </Link>
              <Link href="/sign_up">
                <a
                  className="ml-8 whitespace-nowrap inline-flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-gray-800 hover:bg-white hover:border-gray-800 hover:text-gray-900 transition duration-300"
                >
                  無料登録
                </a>
              </Link>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}
