import React from 'react';
import Link from 'next/link';
import { FaGoogle } from 'react-icons/fa';
import { routes } from '../constants/routes';
import { GOOGLE_LOGIN_URL } from '../constants/api';
import { CurrentAccountProps, withAuth } from '@src/lib/auth';
import { NextPage } from 'next';

const ProfilePage: NextPage<CurrentAccountProps> = ({
  currentAccount,
}: CurrentAccountProps) => {
  return (
    <div className="flex flex-col h-screen bg-gray-50">
      <div className="grid place-items-center mx-2 md:my-20 my-10">
        <div
          className="w-11/12 p-12 sm:w-8/12 md:w-6/12 lg:w-5/12 2xl:w-4/12
            px-6 py-10 sm:px-10 sm:py-6
            bg-white rounded-lg shadow-md lg:shadow-lg"
        >
          <h2 className="text-center font-semibold text-3xl lg:text-4xl text-gray-800">
            Profile
          </h2>

          <div className="text-center m-10">{currentAccount.displayName}</div>
        </div>
      </div>
    </div>
  );
};

export default withAuth(ProfilePage);
