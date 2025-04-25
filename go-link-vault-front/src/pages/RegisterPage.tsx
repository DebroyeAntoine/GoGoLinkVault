import React from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppDispatch, RootState } from "../store";
import { registerUser } from "../features/auth/authSlice";
import { useNavigate } from "react-router-dom";
import AuthForm from "../components/AuthForm";

const RegisterPage = () => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const { loading, error } = useSelector((state: RootState) => state.auth);

  const handleRegister = async (email: string, password: string) => {
    const result = await dispatch(registerUser({ email, password }));
    if (registerUser.fulfilled.match(result)) {
      navigate("/");
    }
  };

  return (
    <AuthForm
      title="Register"
      buttonLabel="Register"
      onSubmit={handleRegister}
      loading={loading}
      error={error}
    />
  );
};

export default RegisterPage;

