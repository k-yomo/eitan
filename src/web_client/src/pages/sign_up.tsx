import React, { FormEvent, ChangeEvent, useState } from 'react';
import { useRouter } from 'next/router';
import Link from 'next/link';
import fetch from 'isomorphic-unfetch';
import { FaGoogle } from 'react-icons/fa';
import { routes } from '../constants/routes';
import {
  CREATE_EMAIL_CONFIRMATION_URL,
  EMAIL_SIGNUP_URL,
  GOOGLE_LOGIN_URL,
} from '../constants/api';

export default function LoginPage() {
  const router = useRouter();
  const [formInput, setFormInput] = useState({
    displayName: '',
    email: '',
    emailConfirmationCode: '',
    password: '',
  });

  const onFormInputChange = (event: ChangeEvent<HTMLInputElement>) => {
    switch (event.target.name) {
      case 'displayName':
        setFormInput({ ...formInput, displayName: event.target.value });
        break;
      case 'email':
        setFormInput({ ...formInput, email: event.target.value });
        break;
      case 'emailConfirmationCode':
        setFormInput({
          ...formInput,
          emailConfirmationCode: event.target.value,
        });
        break;
      case 'password':
        setFormInput({ ...formInput, password: event.target.value });
        break;
    }
  };

  const createEmailConfirmation = () => {
    if (formInput.email === '') {
      return;
    }
    fetch(CREATE_EMAIL_CONFIRMATION_URL, {
      method: 'POST',
      body: JSON.stringify({ email: formInput.email }),
    });
  };

  const onSubmitSignUpForm = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (
      formInput.displayName === '' ||
      formInput.email === '' ||
      formInput.emailConfirmationCode === '' ||
      formInput.password === ''
    ) {
      return;
    }
    fetch(EMAIL_SIGNUP_URL, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      body: JSON.stringify(formInput),
    })
      .then(() => {
        router.push(routes.home());
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
            Sign Up
          </h2>

          <form className="mt-10" method="POST" onSubmit={onSubmitSignUpForm}>
            <label
              htmlFor="displayName"
              className="block text-xs font-semibold text-gray-600 uppercase"
            >
              Display Name
            </label>
            <input
              id="displayName"
              name="displayName"
              placeholder="display name"
              autoComplete="name"
              value={formInput.displayName}
              onChange={onFormInputChange}
              className="block w-full py-3 px-1 mt-2
                    text-gray-800 appearance-none
                    border-b-2 border-gray-100
                    focus:text-gray-500 focus:outline-none focus:border-gray-200"
              required
            />
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
            <div className="flex justify-end">
              <button
                type="button"
                disabled={formInput.email === ''}
                onClick={() => createEmailConfirmation()}
                className="my-1 p-1 rounded-sm text-md
                  whitespace-nowrap border border-transparent shadow-sm text-base font-medium text-white bg-gradient-to-r from-rose-500 to-pink-500 hover:from-rose-600 hover:to-pink-600 transition duration-300 disabled:opacity-50"
              >
                Send confirmation code
              </button>
            </div>

            <label
              htmlFor="emailConfirmationCode"
              className="block mt-2 text-xs font-semibold text-gray-600 uppercase"
            >
              Confirmation Code
            </label>
            <input
              id="emailConfirmationCode"
              name="emailConfirmationCode"
              placeholder="confirmation code"
              value={formInput.emailConfirmationCode}
              onChange={onFormInputChange}
              className="block w-full py-3 px-1 mt-2 mb-4
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
              Sign Up with Email
            </button>
          </form>
          <Link href={GOOGLE_LOGIN_URL}>
            <a>
              <button className="flex flex-row items-center justify-center w-full space-x-2 my-6 p-3 rounded-sm text-md bg-gray-100 hover:bg-gray-200">
                <FaGoogle />
                <div>Sign Up with Google</div>
              </button>
            </a>
          </Link>
          <div className="mt-8 sm:mb-4 text-sm text-right">
            Have an account?{' '}
            <Link href={routes.login()}>
              <a className="flex-2 underline">Log in</a>
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
