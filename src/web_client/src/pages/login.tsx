import React, { ChangeEvent, FormEvent, useState } from 'react';
import Link from 'next/link';
import { FaGoogle } from 'react-icons/fa';
import { routes } from '../constants/routes';
import { EMAIL_LOGIN_URL, GOOGLE_LOGIN_URL } from '../constants/api';
import fetch from 'isomorphic-unfetch';
import { useRouter } from 'next/router';
import { getOriginalUrl, setOriginalUrl } from '@src/lib/originalUrl';

export default function LoginPage() {
  const router = useRouter();
  const [formInput, setFormInput] = useState({
    email: '',
    password: '',
  });

  if (router.query.originalUrl) {
    setOriginalUrl(router.query.originalUrl as string);
  }

  const onFormInputChange = (event: ChangeEvent<HTMLInputElement>) => {
    switch (event.target.name) {
      case 'email':
        setFormInput({ ...formInput, email: event.target.value });
        break;
      case 'password':
        setFormInput({ ...formInput, password: event.target.value });
        break;
    }
  };

  const onSubmitLoginForm = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (formInput.email === '' || formInput.password === '') {
      return;
    }
    fetch(EMAIL_LOGIN_URL, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      body: JSON.stringify(formInput),
    })
      .then(() => {
        router.push(getOriginalUrl() || routes.home());
      })
      .catch((e) => {
        console.log(e);
      });
  };

  return (
    <div className="flex flex-col h-screen bg-gray-50">
      <div className="grid place-items-center mx-2 md:my-20 my-10">
        <div
          className="w-11/12 p-12 sm:w-8/12 md:w-6/12 lg:w-5/12 2xl:w-4/12
            px-6 py-10 sm:px-10 sm:py-6
            bg-white rounded-lg shadow-md lg:shadow-lg"
        >
          <h2 className="text-center font-semibold text-3xl lg:text-4xl text-gray-800">
            Log in
          </h2>

          <form className="mt-10" method="POST" onSubmit={onSubmitLoginForm}>
            <label
              htmlFor="email"
              className="block mt-2 text-xs font-semibold text-gray-600 uppercase"
            >
              Email
            </label>
            <input
              id="email"
              type="email"
              name="email"
              placeholder="email address"
              autoComplete="email"
              value={formInput.email}
              onChange={onFormInputChange}
              className="block w-full py-3 px-1 mt-2
                    text-gray-800 appearance-none
                    border-b-2 border-gray-100
                    focus:text-gray-500 focus:outline-none focus:border-gray-200"
              required
            />

            <label
              htmlFor="password"
              className="block mt-2 text-xs font-semibold text-gray-600 uppercase"
            >
              Password
            </label>
            <input
              id="password"
              type="password"
              name="password"
              placeholder="password"
              autoComplete="current-password"
              value={formInput.password}
              onChange={onFormInputChange}
              className="block w-full py-3 px-1 mt-2 mb-4
                    text-gray-800 appearance-none
                    border-b-2 border-gray-100
                    focus:text-gray-500 focus:outline-none focus:border-gray-200"
              required
            />

            <button
              type="submit"
              className="flex flex-row items-center justify-center w-full space-x-2 mt-10 py-3 rounded-sm text-md
                  whitespace-nowrap border border-transparent shadow-sm text-base font-medium text-white bg-gradient-to-r from-rose-500 to-pink-500 hover:from-rose-600 hover:to-pink-600 transition duration-300"
            >
              Login
            </button>
          </form>
          <Link href={GOOGLE_LOGIN_URL}>
            <a>
              <button className="flex flex-row items-center justify-center w-full space-x-2 my-6 p-3 rounded-sm text-md bg-gray-100 hover:bg-gray-200">
                <FaGoogle />
                <div>Log in with Google</div>
              </button>
            </a>
          </Link>
          <div className="mt-8 sm:mb-4 text-sm text-right">
            Don&apos;t have an account?{' '}
            <Link href={routes.signUp()}>
              <a className="flex-2 underline">Sign Up</a>
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
