import { useAuth } from "@/context/auth";
import LoadingProgress from "../../components/LoadingProgress";
import LoginForm from "../../components/LoginForm";

export default function Auth() {
  const { isLoading } = useAuth();

  if (isLoading) {
    return <LoadingProgress />;
  }

  return <LoginForm />;
}
