import React from 'react';
import Link from 'next/link';
import { Menu, Transition } from '@headlessui/react';
import { useCurrentUserProfileQuery } from '@src/generated/graphql';
import { routes } from '@src/constants/routes';
import { CurrentUserProfileProps } from '@src/lib/auth';
import { LOGOUT_URL } from '@src/constants/api';

export default function Header() {
  return (
    <header>
      <div className="relative bg-white">
        <div className="mx-auto px-4 sm:px-6">
          <div className="flex justify-between items-center border-b-2 border-gray-100 py-4 md:justify-start md:space-x-10">
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
            <RightNav />
          </div>
        </div>
      </div>
    </header>
  );
}

const RightNav = () => {
  const { data, loading, error } = useCurrentUserProfileQuery();
  if (loading) {
    return <></>;
  }

  if (error) {
    console.log(error);
  }

  return (
    <>
      {data && <UserMenu currentUserProfile={data.currentUserProfile} />}
      {(!data || !data.currentUserProfile) && (
        <>
          <div className="md:flex items-center justify-end md:flex-1 lg:w-0">
            <Link href={routes.login()}>
              <a className="whitespace-nowrap text-base font-medium text-gray-500 hover:text-gray-900 transition duration-300">
                Log in
              </a>
            </Link>
            <Link href={routes.signUp()}>
              <a className="ml-8 whitespace-nowrap inline-flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-gray-800 hover:bg-white hover:border-gray-800 hover:text-gray-900 transition duration-300">
                Sign Up
              </a>
            </Link>
          </div>
        </>
      )}
    </>
  );
};

const UserMenu = ({ currentUserProfile }: CurrentUserProfileProps) => {
  return (
    <div className="md:flex items-center justify-end md:flex-1 ">
      <div className="relative inline-block text-left">
        <Menu>
          {({ open }: { open: boolean }) => (
            <>
              <span className="rounded-md shadow-sm">
                <Menu.Button className="w-full px-4 py-2 text-sm font-medium leading-5 text-gray-700 transition duration-150 ease-in-out bg-white rounded-md hover:text-gray-500 focus:outline-none focus:border-blue-300 focus:shadow-outline-blue active:bg-gray-50 active:text-gray-800">
                  {currentUserProfile.screenImgUrl ? (
                    <img
                      className="w-10 h-10 rounded-full"
                      src={currentUserProfile.screenImgUrl}
                      alt={currentUserProfile.displayName}
                    />
                  ) : (
                    <div className="flex justify-center items-center bg-gray-100 w-10 h-10 rounded-full font-bold">
                      {currentUserProfile.displayName.substr(0, 1)}
                    </div>
                  )}
                </Menu.Button>
              </span>

              <Transition
                show={open}
                enter="transition ease-out duration-100"
                enterFrom="transform opacity-0 scale-95"
                enterTo="transform opacity-100 scale-100"
                leave="transition ease-in duration-75"
                leaveFrom="transform opacity-100 scale-100"
                leaveTo="transform opacity-0 scale-95"
              >
                <Menu.Items
                  static
                  className="absolute right-0 w-56 mt-2 origin-top-right bg-white border border-gray-200 divide-y divide-gray-100 rounded-md shadow-lg outline-none"
                >
                  <div className="px-4 py-3">
                    <p className="text-sm leading-5">Logged in as</p>
                    <p className="text-sm font-medium leading-5 text-gray-900 truncate">
                      {currentUserProfile.email}
                    </p>
                  </div>

                  <div className="py-1">
                    <Menu.Item>
                      {({ active }) => (
                        <a
                          href={routes.userSettings()}
                          className={`${
                            active
                              ? 'bg-gray-100 text-gray-900'
                              : 'text-gray-700'
                          } flex justify-between w-full px-4 py-2 text-sm leading-5 text-left`}
                        >
                          User settings
                        </a>
                      )}
                    </Menu.Item>
                  </div>

                  <div className="py-1">
                    <Menu.Item>
                      {({ active }) => (
                        <a
                          href={LOGOUT_URL}
                          className={`${
                            active
                              ? 'bg-gray-100 text-gray-900'
                              : 'text-gray-700'
                          } flex justify-between w-full px-4 py-2 text-sm leading-5 text-left`}
                        >
                          Log out
                        </a>
                      )}
                    </Menu.Item>
                  </div>
                </Menu.Items>
              </Transition>
            </>
          )}
        </Menu>
      </div>
    </div>
  );
};
