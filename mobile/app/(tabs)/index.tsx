import LoginForm from "@/components/LoginForm";
import OutingsPage from "@/components/OutingsPage";
import { useAuth } from "@/context/auth";
import LoadingProgress from "@/components/LoadingProgress";

export default function HomeScreen() {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return <LoadingProgress />;
  }

  if (!isLoading && !user) {
    return <LoginForm />;
  }

  return <OutingsPage />;
}
