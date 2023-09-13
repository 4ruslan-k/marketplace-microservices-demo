"use client";

import React, { useContext } from 'react';
import TextField from '@mui/material/TextField';
import * as yup from 'yup';
import { useFormik } from 'formik';
import { useSearchParams,useRouter } from 'next/navigation';
import { emailValidation, passwordValidation, nameValidation } from '../constants/validations';
import AuthCardWrapper from '../components/AuthCardWrapper';
import { SnackbarContext } from '../context/SnackbarContext';
import { useFetchUser, useSignupUser } from '../requests/userHooks';

export default function SignUp() {
  const { openSuccessSnackbar, openErrorSnackbar } = useContext(SnackbarContext);
  const { refetchUser } = useFetchUser();

  const router = useRouter();

  const { signupUser, isLoadingSignup } = useSignupUser({
    onSuccess: () => {
      openSuccessSnackbar('User created successfully');
      refetchUser();
      router.push('/');
    },
    onError: (error) => {
      openErrorSnackbar(error?.response?.data?.message || 'Something went wrong');
    },
  });
  const signUpSchema = yup.object().shape({
    name: nameValidation,
    email: emailValidation,
    password: passwordValidation,
    confirmPassword: yup
      .string()
      .required('This field is required.')
      .oneOf([yup.ref('password')], 'Passwords not match'),
  });

  const { errors, handleChange, touched, handleBlur, isValid, handleSubmit } = useFormik({
    initialValues: {
      name: '',
      email: '',
      password: '',
      confirmPassword: '',
    },
    validationSchema: signUpSchema,
    onSubmit: signupUser,
  });

  const nameError = errors.name && touched.name ? errors.name : null;
  const emailError = errors.email && touched.email ? errors.email : null;
  const passwordError = errors.password && touched.password ? errors.password : null;
  const confirmPasswordError = errors.confirmPassword && touched.confirmPassword ? errors.confirmPassword : null;
  return (
    <AuthCardWrapper
      isValid={isValid}
      handleSubmit={handleSubmit}
      title="Create An Account"
      submitButtonLabel="Sign Up"
      isLoading={isLoadingSignup}
    >
      <div className="mb-10">
        <TextField
          className="w-full"
          id="name"
          name="name"
          label="Name"
          variant="outlined"
          onChange={handleChange}
          onBlur={handleBlur}
          helperText={nameError}
          error={Boolean(nameError)}
          required
        />
      </div>
      <div className="mb-10">
        <TextField
          className="w-full"
          id="email"
          name="email"
          label="Email"
          type="email"
          variant="outlined"
          onChange={handleChange}
          onBlur={handleBlur}
          helperText={emailError}
          error={Boolean(emailError)}
          required
        />
      </div>
      <div className="mb-10">
        <TextField
          className="w-full"
          id="password"
          name="password"
          label="Password"
          type="password"
          variant="outlined"
          onChange={handleChange}
          onBlur={handleBlur}
          helperText={passwordError}
          error={Boolean(passwordError)}
          required
        />
      </div>
      <TextField
        className="w-full"
        id="confirmPassword"
        name="confirmPassword"
        label="Confirm Password"
        type="password"
        variant="outlined"
        onChange={handleChange}
        onBlur={handleBlur}
        helperText={confirmPasswordError}
        error={Boolean(confirmPasswordError)}
        required
      />
    </AuthCardWrapper>
  );
}
